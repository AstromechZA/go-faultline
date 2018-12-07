# go-faultline

![Golang](https://img.shields.io/badge/language-golang-green.svg?style=flat-square)
![License](https://img.shields.io/badge/license-MIT-green.svg?style=flat-square)
![Dependencies](https://img.shields.io/badge/dependencies-0-green.svg?style=flat-square)

Fault injection library and api's for Golang.

**NOTE**: this is just an idea so far and is still in development and prototyping.

See [examples/server.go](examples/server.go) for an example..

## Why?

This is an experimental form of chaos engineering for internal API's that can be used to inject failures with a certain 
probability. For example: cause 0.1% of requests to have additional latency or to return a 500 error.

By having tunable probablities and delays, and by outputting an audit log of the effects, this fault injection library
adds a permanent level of unreliability to your API while helping you to debug when faults have been injected. This sounds bad, but remember that this can be used to enforce an 
error budget or SLO and will help to force consumers of the API to gracefully handle these faults. 

The idea is that in development or pre-production you can run with the probabilities raised, and when you are 
in production you can turn it off or selectively lower it.
