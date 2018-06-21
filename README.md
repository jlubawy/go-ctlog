# jlubawy/go-ctlog

[![GoDoc](https://godoc.org/github.com/jlubawy/go-ctlog?status.svg)](https://godoc.org/github.com/jlubawy/go-ctlog)
[![Build Status](https://travis-ci.org/jlubawy/go-ctlog.svg?branch=master)](https://travis-ci.org/jlubawy/go-ctlog)

This project is a C library and a collection of build tools that can be used to
add tokenized logging to a C project.

_This is a work in progress, no guarantees are made as far as API goes. If there
is interest in this project let me know and I can work towards making more
guarantees._

The idea behind tokenized logging is to reduce program binary sizes by replacing
strings used for debugging and logging with "tokens". A token is defined as a
module index and line number that can be used to lookup the replaced strings at
runtime to reproduce the original output. __In real-life programs the savings
are typically multiple kB which can make or break a project if memory is at a
premium.__

See [examples/arduino](examples/arduino) for a more concrete example running on
an embedded device.

See [examples/basic](examples/basic) for something that can be run on a PC.
