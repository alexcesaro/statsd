# statsd

## Introduction

statsd is a simple and efficient [Statsd](https://github.com/etsy/statsd)
client.

This is a fork of a [dead statsd client project](https://github.com/alexcesaro/statsd) with various impprovements.

## Features

- Supports all StatsD metrics: counter, gauge (absolute and relative), timing and set
- Supports InfluxDB and Datadog tags
- Fast and GC-friendly: all functions for sending metrics do not allocate
- Efficient: metrics are buffered by default
- Simple and clean API

## Download

    go get github.com/multiplay/statsd

## License

[MIT](LICENSE)


## Contribute

Do you have any question the documentation does not answer? Is there a use case
that you feel is common and is not well-addressed by the current API?

If so you are more than welcome to ask questions or open an issue or send a pull-request here on Github.
