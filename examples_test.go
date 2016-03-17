package statsd_test

import (
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/alexcesaro/statsd"
)

func Example() {
	// Connect to the UDP port 8125 by default.
	c, err := statsd.New()
	if err != nil {
		// If nothing is listening on the target port, an error is returned and
		// the returned client does nothing but is still usable. So we can
		// just log the error and go on.
		log.Print(err)
	}
	defer c.Close()

	// Increment a counter.
	c.Increment("foo.counter")

	// Gauge something.
	c.Gauge("num_goroutine", runtime.NumGoroutine())

	// Time something.
	t := c.NewTiming()
	http.Get("http://example.com/")
	t.Send("homepage.response_time")

	// It can also be used as a one-liner to easily time a function.
	pingHomepage := func() {
		defer c.NewTiming().Send("homepage.response_time")

		http.Get("http://example.com/")
	}
	pingHomepage()
}

func ExampleClient_Clone() {
	c, err := statsd.New(statsd.Prefix("my_app"))
	if err != nil {
		log.Print(err)
	}

	httpStats := c.Clone(statsd.Prefix("http"))
	httpStats.Increment("foo.bar") // Increments: my_app.http.foo.bar
}

func ExampleAddress() {
	statsd.New(statsd.Address("192.168.0.5:8126"))
}

func ExampleErrorHandler() {
	statsd.New(statsd.ErrorHandler(func(err error) {
		log.Print(err)
	}))
}

func ExampleFlushPeriod() {
	statsd.New(statsd.FlushPeriod(10 * time.Millisecond))
}

func ExampleMaxPacketSize() {
	statsd.New(statsd.MaxPacketSize(512))
}

func ExampleNetwork() {
	// Send metrics using a TCP connection.
	statsd.New(statsd.Network("tcp"))
}

func ExampleTagsFormat() {
	statsd.New(statsd.TagsFormat(statsd.InfluxDB))
}

func ExampleMute() {
	c, err := statsd.New(statsd.Mute(true))
	if err != nil {
		log.Print(err)
	}
	c.Increment("foo.bar") // Does nothing.
}

func ExamplePrefix() {
	c, err := statsd.New(statsd.Prefix("my_app"))
	if err != nil {
		log.Print(err)
	}
	c.Increment("foo.bar") // Increments: my_app.foo.bar
}

func ExampleSampleRate() {
	statsd.New(statsd.SampleRate(0.2)) // Sends metrics 20% of the time.
}

func ExampleTags() {
	statsd.New(
		statsd.TagsFormat(statsd.InfluxDB),
		statsd.Tags("region", "us", "app", "my_app"),
	)
}

var c *statsd.Client

func ExampleClient_NewTiming() {
	// Send a timing metric each time the function is run.
	defer c.NewTiming().Send("homepage.response_time")
	http.Get("http://example.com/")
}
