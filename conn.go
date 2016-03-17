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
	mu sync.Mutex

	// Fields guarded by the mutex.
	w         io.WriteCloser
	buf       []byte
	rateCache map[float32]string
	closed    bool

	// Fields settable with options at Client's creation.
	addr          string
	errorHandler  func(error)
	flushPeriod   time.Duration
	maxPacketSize int
	network       string
	tagFormat     TagFormat
}

func newConn(conf connConfig, muted bool) (*conn, error) {
	c := &conn{
		addr:          conf.Addr,
		errorHandler:  conf.ErrorHandler,
		flushPeriod:   conf.FlushPeriod,
		maxPacketSize: conf.MaxPacketSize,
		network:       conf.Network,
		tagFormat:     conf.TagFormat,
	}

	if muted {
		return c, nil
	}

	var err error
	c.w, err = dialTimeout(c.network, c.addr, 5*time.Second)
	if err != nil {
		return c, err
	}
	// When using UDP do a quick check to see if something is listening on the
	// given port to return an error as soon as possible.
	if c.network[:3] == "udp" {
		for i := 0; i < 2; i++ {
			_, err = c.w.Write(nil)
			if err != nil {
				_ = c.w.Close()
				c.w = nil
				return c, err
			}
		}
	}

	// To prevent a buffer overflow add some capacity to the buffer to allow for
	// an additional metric.
	c.buf = make([]byte, 0, c.maxPacketSize+200)

	if c.flushPeriod > 0 {
		go func() {
			ticker := time.NewTicker(c.flushPeriod)
			for _ = range ticker.C {
				c.mu.Lock()
				if c.closed {
					ticker.Stop()
					c.mu.Unlock()
					return
				}
				c.flush(0)
				c.mu.Unlock()
			}
		}()
	}

	return c, nil
}

func (c *conn) count(prefix, bucket string, n int, rate float32, tags string) {
	c.mu.Lock()
	l := len(c.buf)
	c.appendBucket(prefix, bucket, tags)
	c.appendInt(n)
	c.appendType("c")
	c.appendRate(rate)
	c.closeMetric(tags)
	c.flushIfBufferFull(l)
	c.mu.Unlock()
}

func (c *conn) gauge(prefix, bucket string, value int, tags string) {
	c.mu.Lock()
	l := len(c.buf)
	// To set a gauge to a negative value we must first set it to 0.
	// https://github.com/etsy/statsd/blob/master/docs/metric_types.md#gauges
	if value < 0 {
		c.appendBucket(prefix, bucket, tags)
		c.appendGauge(0, tags)
	}
	c.appendBucket(prefix, bucket, tags)
	c.appendGauge(value, tags)
	c.flushIfBufferFull(l)
	c.mu.Unlock()
}

func (c *conn) appendGauge(value int, tags string) {
	c.appendInt(value)
	c.appendType("g")
	c.closeMetric(tags)
}

func (c *conn) timing(prefix, bucket string, n int, rate float32, tags string) {
	c.mu.Lock()
	l := len(c.buf)
	c.appendBucket(prefix, bucket, tags)
	c.appendInt(n)
	c.appendType("ms")
	c.appendRate(rate)
	c.closeMetric(tags)
	c.flushIfBufferFull(l)
	c.mu.Unlock()
}

func (c *conn) histogram(prefix, bucket string, n int, rate float32, tags string) {
	c.mu.Lock()
	l := len(c.buf)
	c.appendBucket(prefix, bucket, tags)
	c.appendInt(n)
	c.appendType("h")
	c.appendRate(rate)
	c.closeMetric(tags)
	c.flushIfBufferFull(l)
	c.mu.Unlock()
}

func (c *conn) unique(prefix, bucket string, value string, tags string) {
	c.mu.Lock()
	l := len(c.buf)
	c.appendBucket(prefix, bucket, tags)
	c.appendString(value)
	c.appendType("s")
	c.closeMetric(tags)
	c.flushIfBufferFull(l)
	c.mu.Unlock()
}

func (c *conn) appendByte(b byte) {
	c.buf = append(c.buf, b)
}

func (c *conn) appendString(s string) {
	c.buf = append(c.buf, s...)
}

func (c *conn) appendInt(i int) {
	c.buf = strconv.AppendInt(c.buf, int64(i), 10)
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

func (c *conn) flushIfBufferFull(lastSafeLen int) {
	if len(c.buf) > c.maxPacketSize {
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

	// Trim the last \n, StatsD does not like it.
	_, err := c.w.Write(c.buf[:n-1])
	c.handleError(err)
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
