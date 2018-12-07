# go-faultline

Fault injection library and api's for Golang.

**NOTE**: this is just an idea so far and is still in development and prototyping.

See [examples/server.go](examples/server.go) for an example..

## Why?

This is an experimental form of chaos engineering for internal API's that can be used to inject failures with a certain 
probability. For example: cause 0.1% of requests to have additional latency or to return a 500 error.

By having tunable probablities and delays, and by outputting an audit log of the effects, this fault injection library
adds a permanent level of unreliability to your API while helping you to debug when faults have been injected. This sounds bad, but remember that this can be used to enforce an 
error budget or SLO and will help to force consumers of the API to gracefully handle these faults. 
