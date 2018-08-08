# GoLang jrpc2Client (early beta)

[Website](https://riftbit.com) | [Blog](https://ergoz.ru/)

[![license](https://img.shields.io/github/license/riftbit/jrpc2Client.svg)](LICENSE)
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/riftbit/jrpc2Client)
[![Coverage Status](https://coveralls.io/repos/github/riftbit/jrpc2Client/badge.svg?branch=master)](https://coveralls.io/github/riftbit/jrpc2Client?branch=master)
[![Build Status](https://travis-ci.org/riftbit/jrpc2Client.svg?branch=master)](https://travis-ci.org/riftbit/jrpc2Client)
[![Go Report Card](https://goreportcard.com/badge/github.com/riftbit/jrpc2Client)](https://goreportcard.com/report/github.com/riftbit/jrpc2Client)

This is a json-rpc 2.0 client package for golang based on:

 - **HTTP Client:** [erikdubbelboer/fasthttp](github.com/erikdubbelboer/fasthttp)
 - **JSON Parser:** [pquerna/ffjson](github.com/pquerna/ffjson/ffjson)
 - **Logger:** [Sirupsen/logrus](github.com/Sirupsen/logrus)

to get high perfomance

This package is still in development

## Examples

### Without custom logger settings

```golang
package main

import (
	"github.com/riftbit/jrpc2Client"
)

type TestReply struct {
	LogID string `json:"log_id"`
}

func main() {
	client := jrpc2Client.NewClient()

	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuth("user", "password")

	dstT := &TestReply{}
	err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_1"}, dstT)
	if err != nil {
		panic(err)
	}
	println(dstT.LogID)
}
```


### With custom logger settings

```golang
package main

import (
	"github.com/riftbit/jrpc2Client"
)

type TestReply struct {
	LogID string `json:"log_id"`
}

func main() {
	logger := &logrus.Logger{
    		Out:       os.Stdout,
    		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
    		Level:     logrus.DebugLevel,
    }

    client := jrpc2Client.NewClientWithLogger(logger)

    client.SetBaseURL("http://127.0.0.1:65001")
    client.SetUserAgent("JsonRPC Test Client")
    client.SetBasicAuth("user", "password")

    dstT := &TestReply{}
    err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_3"}, dstT)
    if err != nil {
    		panic(err)
    }
    println(dstT.LogID)
}
```


## Benchmark results
