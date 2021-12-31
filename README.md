This is a backwards compatible fork of [alexcesaro/statsd](https://github.com/alexcesaro/statsd).

Note that this repository is versioned independently of the (unmaintained) upstream.

**Changelog**

* 2022-01-01 - Added support for relative gauges via `Client.GaugeRelative`, and fixed gauge handling of -0 floats

* 2021-12-31 - Added support for disabling the initial ping/check for UDP connections via new `UDPCheck` option

* 2020-11-13 - Added support for go modules

* 2020-11-13 - Added the SafeConn write closer that checks if the connection is still up before attempting to write

* 2020-08-27 - Fixed bug in `Tags` identified by https://github.com/alexcesaro/statsd/issues/41

* 2019-05-22 - Added support for arbitrary output streams via new `WriteCloser` option

* 2019-05-22 - Added support for simplified inline flush logic via new `InlineFlush` option

* 2019-05-26 - Fixed bug causing trailing newlines to be removed for streaming (non-udp) connections

See also the [upstream changelog](CHANGELOG.md).

---

# statsd
[![Code Coverage](https://gocover.io/_badge/github.com/joeycumines/statsd)](https://gocover.io/github.com/joeycumines/statsd)
[![Documentation](https://godoc.org/github.com/joeycumines/statsd?status.svg)](https://godoc.org/github.com/joeycumines/statsd)

## Introduction

Package statsd is a simple and efficient [StatsD](https://github.com/statsd/statsd)
client.

See the [benchmark](https://github.com/alexcesaro/statsdbench) for a comparison
with other Go StatsD clients.

You may also be interested in this (old / historic)
[thread on golang-nuts](https://groups.google.com/d/topic/golang-nuts/Tz6t4_iLgnw/discussion).

## Features

- Supports all StatsD metrics: counter, gauge (absolute and relative), timing and set
- Supports InfluxDB and Datadog tags
- Fast and GC-friendly: all functions for sending metrics do not allocate
- Efficient: metrics are buffered by default
- Simple and clean API
- 100% test coverage
- Versioned API using gopkg.in


## Documentation

https://pkg.go.dev/github.com/joeycumines/statsd#section-documentation


## Download

    go get github.com/joeycumines/statsd


## Example

See the [examples in the documentation](https://pkg.go.dev/github.com/joeycumines/statsd#pkg-examples).


## License

[MIT](LICENSE)


## Contribute

Do you have any question the documentation does not answer? Is there a use case
that you feel is common and is not well-addressed by the current API?

If so you are more than welcome to ask questions, open an issue, or send a
pull request here on GitHub.
