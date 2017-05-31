# statsd
[![Build Status](https://travis-ci.org/alexcesaro/statsd.svg?branch=v2)](https://travis-ci.org/alexcesaro/statsd) [![Code Coverage](http://gocover.io/_badge/github.com/SurfEasy/statsd-1)](http://gocover.io/github.com/SurfEasy/statsd-1) [![Documentation](https://godoc.org/github.com/SurfEasy/statsd-1?status.svg)](https://godoc.org/github.com/SurfEasy/statsd-1)

## Introduction

statsd is a simple and efficient [Statsd](https://github.com/etsy/statsd)
client.

See the [benchmark](https://github.com/alexcesaro/statsdbench) for a comparison
with other Go StatsD clients.

## Features

- Supports all StatsD metrics: counter, gauge, timing and set
- Supports InfluxDB and Datadog tags
- Fast and GC-friendly: all functions for sending metrics do not allocate
- Efficient: metrics are buffered by default
- Simple and clean API
- 100% test coverage
- Versioned API using gopkg.in


## Documentation

https://godoc.org/github.com/SurfEasy/statsd-1


## Download

    go get github.com/SurfEasy/statsd-1


## Example

See the [examples in the documentation](https://godoc.org/github.com/SurfEasy/statsd-1#example-package).


## License

[MIT](LICENSE)


## Contribute

Do you have any question the documentation does not answer? Is there a use case
that you feel is common and is not well-addressed by the current API?

If so you are more than welcome to ask questions in the
[thread on golang-nuts](https://groups.google.com/d/topic/golang-nuts/Tz6t4_iLgnw/discussion)
or open an issue or send a pull-request here on Github.
