# statsd
[![Build Status](https://travis-ci.org/alexcesaro/statsd.svg?branch=v1)](https://travis-ci.org/alexcesaro/statsd) [![Code Coverage](http://gocover.io/_badge/gopkg.in/alexcesaro/statsd.v1)](http://gocover.io/gopkg.in/alexcesaro/statsd.v1) [![Documentation](https://godoc.org/gopkg.in/alexcesaro/statsd.v1?status.svg)](https://godoc.org/gopkg.in/alexcesaro/statsd.v1)

## Introduction

statsd is a simple and efficient [Statsd](https://github.com/etsy/statsd)
client.

See the [benchmark](https://github.com/alexcesaro/statsdbench) for a comparison
with other Go StatsD clients.

## Features

- Supports all StatsD metrics: counter, gauge, timing and set
- Fast and GC-friendly: Client's methods do not allocate
- Simple API
- 100% test coverage
- Versioned API using gopkg.in


## Documentation

https://godoc.org/gopkg.in/alexcesaro/statsd.v1


## Download

    go get gopkg.in/alexcesaro/statsd.v1


## Example

See the [examples in the documentation](https://godoc.org/gopkg.in/alexcesaro/statsd.v1#example-package).


## License

[MIT](LICENSE)
