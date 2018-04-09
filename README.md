# statsd

## Introduction

statsd is a simple and efficient [Statsd](https://github.com/etsy/statsd)
client.

This is a fork of the [Go Statsd client](https://github.com/alexcesaro/statsd) by [alexcesaro](https://github.com/alexcesaro).
This features some additions to gauges and listener-less client startup.

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

    go get github.com/infosum/statsd


## Example

See the [examples in the documentation](https://godoc.org/gopkg.in/alexcesaro/statsd.v2#example-package).


## License

[MIT](LICENSE)


## Contribute

Do you have any question the documentation does not answer? Is there a use case
that you feel is common and is not well-addressed by the current API?

If so you are more than welcome to ask questions or open an issue or send a pull-request here on Github.
