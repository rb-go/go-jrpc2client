package jrpc2client_test

import (
	"fmt"
	"time"

	"github.com/riftbit/jrpc2client"
)

// This function is named ExampleClient_Call(), this way godoc knows to associate
// it with the Client type and method Call.
func ExampleClient_Call() {

	type TestReply struct {
		LogID     string `json:"log_id"`
		UserAgent string `json:"user_agent"`
	}

	type TestArgs struct {
		ID string
	}

	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:8080")
	dstP := &TestReply{}
	client.SetClientTimeout(10 * time.Millisecond)

	// final url will be http://127.0.0.1:8080/api
	err := client.Call("/api", "demo.TestClientTimeout", TestArgs{ID: "123"}, dstP)
	if err != nil {
		panic(err)
	}

	fmt.Println(dstP.LogID)
}

// This function is named ExampleClient_CallForMap(), this way godoc knows to associate
// it with the Client type and method CallForMap.
func ExampleClient_CallForMap() {
	type TestArgs struct {
		ID string
	}

	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:8080")

	client.SetClientTimeout(10 * time.Millisecond)
	// final url will be http://127.0.0.1:8080/api
	dstM, err := client.CallForMap("/api", "demo.TestClientTimeout", TestArgs{ID: "321"})
	if err != nil {
		panic(err)
	}

	val, ok := dstM["log_id"]
	if ok != true {
		panic("key log_id not exists")
	}

	fmt.Println(val)
}

func ExampleNewClient() {
	client := jrpc2client.NewClient()
	//print empty line because BaseURL not setted
	println(client.BaseURL)
}
