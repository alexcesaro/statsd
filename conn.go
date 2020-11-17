package statsd

import (
	"io"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

type conn struct {
	// config

	errorHandler  func(error)
	flushPeriod   time.Duration
	maxPacketSize int
	tagFormat     TagFormat
	inlineFlush   bool

	// state

	mu                  sync.Mutex         // mu synchronises internal state
	closed              bool               // closed indicates if w has been closed (triggered by first client close)
	w                   io.WriteCloser     // w is the writer for the connection
	buf                 []byte             // buf is the buffer for the connection
	rateCache           map[float32]string // rateCache caches string representations of sampling rates
	trimTrailingNewline bool               // trimTrailingNewline is set only when running in UDP mode
}

func newConn(conf connConfig, muted bool) (*conn, error) {
	c := &conn{
		errorHandler:  conf.ErrorHandler,
		flushPeriod:   conf.FlushPeriod,
		maxPacketSize: conf.MaxPacketSize,
		tagFormat:     conf.TagFormat,
		inlineFlush:   conf.InlineFlush,
		w:             conf.WriteCloser,
	}

	// exit if muted
	if muted {
		// close and clear any provided writer
		if c.w != nil {
			_ = c.w.Close()
			c.w = nil
		}
		// return muted client
		return c, nil
	}

	// initialise writer if not provided
	if c.w == nil {
		if err := c.connect(conf.Network, conf.Addr); err != nil {
			return c, err
		}
	}

	// To prevent a buffer overflow add some capacity to the buffer to allow for
	// an additional metric.
	c.buf = make([]byte, 0, c.maxPacketSize+200)

	// start the flush worker only if we have a rate and it's not unnecessary
	if c.flushPeriod > 0 && !c.inlineFlush {
		go c.flushWorker()
	}

	return c, nil
}

func (c *conn) flushWorker() {
	ticker := time.NewTicker(c.flushPeriod)
	defer ticker.Stop()
	for range ticker.C {
		if func() bool {
			c.mu.Lock()
			defer c.mu.Unlock()
			if c.closed {
				return true
			}
			c.flush(0)
			return false
		}() {
			return
		}
	}
}

func (c *conn) connect(network string, address string) error {
	var err error
	c.w, err = dialTimeout(network, address, 5*time.Second)
	if err != nil {
		return err
	}

	if network[:3] == "udp" {
		// udp retains behavior from the original implementation where it would strip a trailing newline
		c.trimTrailingNewline = true

		// When using UDP do a quick check to see if something is listening on the
		// given port to return an error as soon as possible.
		for i := 0; i < 2; i++ {
			_, err = c.w.Write(nil)
			if err != nil {
				_ = c.w.Close()
				c.w = nil
				return err
			}
		}
	}
	return nil
}

func (c *conn) metric(prefix, bucket string, n interface{}, typ string, rate float32, tags string) {
	c.mu.Lock()
	l := len(c.buf)
	c.appendBucket(prefix, bucket, tags)
	c.appendNumber(n)
	c.appendType(typ)
	c.appendRate(rate)
	c.closeMetric(tags)
	c.flushIfNecessary(l)
	c.mu.Unlock()
}

func (c *conn) gauge(prefix, bucket string, value interface{}, tags string) {
	c.mu.Lock()
	l := len(c.buf)
	// To set a gauge to a negative value we must first set it to 0.
	// https://github.com/etsy/statsd/blob/master/docs/metric_types.md#gauges
	if isNegative(value) {
		c.appendBucket(prefix, bucket, tags)
		c.appendGauge(0, tags)
	}
	c.appendBucket(prefix, bucket, tags)
	c.appendGauge(value, tags)
	c.flushIfNecessary(l)
	c.mu.Unlock()
}

func (c *conn) appendGauge(value interface{}, tags string) {
	c.appendNumber(value)
	c.appendType("g")
	c.closeMetric(tags)
}

func (c *conn) unique(prefix, bucket string, value string, tags string) {
	c.mu.Lock()
	l := len(c.buf)
	c.appendBucket(prefix, bucket, tags)
	c.appendString(value)
	c.appendType("s")
	c.closeMetric(tags)
	c.flushIfNecessary(l)
	c.mu.Unlock()
}

func (c *conn) appendByte(b byte) {
	c.buf = append(c.buf, b)
}

func (c *conn) appendString(s string) {
	c.buf = append(c.buf, s...)
}

func (c *conn) appendNumber(v interface{}) {
	switch n := v.(type) {
	case int:
		c.buf = strconv.AppendInt(c.buf, int64(n), 10)
	case uint:
		c.buf = strconv.AppendUint(c.buf, uint64(n), 10)
	case int64:
		c.buf = strconv.AppendInt(c.buf, n, 10)
	case uint64:
		c.buf = strconv.AppendUint(c.buf, n, 10)
	case int32:
		c.buf = strconv.AppendInt(c.buf, int64(n), 10)
	case uint32:
		c.buf = strconv.AppendUint(c.buf, uint64(n), 10)
	case int16:
		c.buf = strconv.AppendInt(c.buf, int64(n), 10)
	case uint16:
		c.buf = strconv.AppendUint(c.buf, uint64(n), 10)
	case int8:
		c.buf = strconv.AppendInt(c.buf, int64(n), 10)
	case uint8:
		c.buf = strconv.AppendUint(c.buf, uint64(n), 10)
	case float64:
		c.buf = strconv.AppendFloat(c.buf, n, 'f', -1, 64)
	case float32:
		c.buf = strconv.AppendFloat(c.buf, float64(n), 'f', -1, 32)
	}
}

func isNegative(v interface{}) bool {
	switch n := v.(type) {
	case int:
		return n < 0
	case uint:
		return n < 0
	case int64:
		return n < 0
	case uint64:
		return n < 0
	case int32:
		return n < 0
	case uint32:
		return n < 0
	case int16:
		return n < 0
	case uint16:
		return n < 0
	case int8:
		return n < 0
	case uint8:
		return n < 0
	case float64:
		return n < 0
	case float32:
		return n < 0
	}
	return false
}

func (c *conn) appendBucket(prefix, bucket string, tags string) {
	c.appendString(prefix)
	c.appendString(bucket)
	if c.tagFormat == InfluxDB {
		c.appendString(tags)
	}
	c.appendByte(':')
}

func (c *conn) appendType(t string) {
	c.appendByte('|')
	c.appendString(t)
}

func (c *conn) appendRate(rate float32) {
	if rate == 1 {
		return
	}
	if c.rateCache == nil {
		c.rateCache = make(map[float32]string)
	}

	c.appendString("|@")
	if s, ok := c.rateCache[rate]; ok {
		c.appendString(s)
	} else {
		s = strconv.FormatFloat(float64(rate), 'f', -1, 32)
		c.rateCache[rate] = s
		c.appendString(s)
	}
}

func (c *conn) closeMetric(tags string) {
	if c.tagFormat == Datadog {
		c.appendString(tags)
	}
	c.appendByte('\n')
}

func (c *conn) flushNecessary() bool {
	if c.inlineFlush {
		return true
	}
	if len(c.buf) > c.maxPacketSize {
		return true
	}
	return false
}

func (c *conn) flushIfNecessary(lastSafeLen int) {
	if c.inlineFlush {
		lastSafeLen = 0
	}
	if c.flushNecessary() {
		c.flush(lastSafeLen)
	}
}

// flush flushes the first n bytes of the buffer.
// If n is 0, the whole buffer is flushed.
func (c *conn) flush(n int) {
	if len(c.buf) == 0 {
		return
	}
	if n == 0 {
		n = len(c.buf)
	}

	// write
	buffer := c.buf[:n]
	if c.trimTrailingNewline {
		// https://github.com/cactus/go-statsd-client/issues/17
		// Trim the last \n, StatsD does not like it.
		buffer = buffer[:len(buffer)-1]
	}
	_, err := c.w.Write(buffer)
	c.handleError(err)

	// consume
	if n < len(c.buf) {
		copy(c.buf, c.buf[n:])
	}
	c.buf = c.buf[:len(c.buf)-n]
}

func (c *conn) handleError(err error) {
	if err != nil && c.errorHandler != nil {
		c.errorHandler(err)
	}
}

// Stubbed out for testing.
var (
	dialTimeout = net.DialTimeout
	now         = time.Now
	randFloat   = rand.Float32
)
