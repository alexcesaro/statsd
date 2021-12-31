This is a backwards compatible fork of [alexcesaro/statsd](https://github.com/alexcesaro/statsd).

Note that this repository is versioned independently of the (unmaintained) upstream.

**Changelog**

* 2020-11-13 - Added support for go modules

* 2020-11-13 - Added the SafeConn write closer that checks if the connection is still up before attempting to write

* 2020-08-27 - Fixed bug in `statsd.Tags` identified by https://github.com/alexcesaro/statsd/issues/41

* 2019-05-22 - Added support for arbitrary output streams via new `statsd.WriteCloser` option

* 2019-05-22 - Added support for simplified inline flush logic via new `statsd.InlineFlush` option

* 2019-05-26 - Fixed bug causing trailing newlines to be removed for streaming (non-udp) connections

See also the [upstream changelog](CHANGELOG.md).

---

# statsd
[![Build Status](https://travis-ci.org/alexcesaro/statsd.svg?branch=v2)](https://travis-ci.org/alexcesaro/statsd) [![Code Coverage](http://gocover.io/_badge/gopkg.in/alexcesaro/statsd.v2)](http://gocover.io/gopkg.in/alexcesaro/statsd.v2) [![Documentation](https://godoc.org/gopkg.in/alexcesaro/statsd.v2?status.svg)](https://godoc.org/gopkg.in/alexcesaro/statsd.v2)

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

https://godoc.org/gopkg.in/alexcesaro/statsd.v2


## Download

    go get github.com/joeycumines/statsd


## Example

See the [examples in the documentation](https://godoc.org/gopkg.in/alexcesaro/statsd.v2#example-package).


## License

[MIT](LICENSE)


## Contribute

Do you have any question the documentation does not answer? Is there a use case
that you feel is common and is not well-addressed by the current API?

If so you are more than welcome to ask questions in the
[thread on golang-nuts](https://groups.google.com/d/topic/golang-nuts/Tz6t4_iLgnw/discussion)
or open an issue or send a pull-request here on Github.
