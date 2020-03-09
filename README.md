# go-gcache-redis-sample
This project is for experimenting with using
[gcache](https://github.com/bluele/gcache) as a loading cache that's ultimately
backed by [Redis](https://redis.io/).

The application instance will keep a smaller instance memory backed cache that
looks up a key from Redis if the gcache doesn't have the value.

## Getting started
A `Makefile` is available with convenience commands:
```
⇒  make help
help                           View help information
dev                            Runs a local development environment
```
