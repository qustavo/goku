goku [![Build Status](https://travis-ci.org/gchaincl/goku.svg?branch=master)](https://travis-ci.org/gchaincl/goku) [![GoDoc](https://godoc.org/github.com/gchaincl/goku?status.svg)](https://godoc.org/github.com/gchaincl/goku)
===

Idiomatic queues for different providers.  
Goku tries to unify queue interaction (currently only SQS is supported only),
by using channels.

Usage
---

First, you need to initialize the library, specifying a provider, in this case
`sqs_provider` is a [goamz/SQS](http://godoc.org/github.com/crowdmob/goamz/sqs) instance.

```go
q := goku.NewQueue(
	goku.Reader(sqs_provider),
	goku.Writer(sqs_provider),
)
```

`goku.Queue` will expose two channels, one for reading messages `Sender()`,
and one for writing messages `Receiver()`.

After that, you should be able to exchange messages with the queue as follow:
```go
q.Receiver() <- "Message"
msg := <-q.Sender()
```

For more information, please read the [example](https://github.com/gchaincl/goku/blob/master/examples/sqs/sqs.go).
