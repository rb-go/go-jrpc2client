package jrpc2client

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/riftbit/jrpc2server"
	"github.com/stretchr/testify/assert"
)

// DemoAPI area
type DemoMapAPI struct{}

type TestMapArgs struct {
	ID string
}

type TestMapReply struct {
	LogID string `json:"log_id"`
}

// Test Method to test
func (h *DemoMapAPI) Test(ctx *fasthttp.RequestCtx, args *TestMapArgs, reply *TestMapReply) error {
	if args.ID == "" {
		return &jrpc2server.Error{Code: jrpc2server.JErrorInvalidParams, Message: "ID should not be empty"}
	}
	reply.LogID = args.ID
	return nil
}

func TestMapPrepare(t *testing.T) {
	api := jrpc2server.NewServer()
	err := api.RegisterService(new(DemoMapAPI), "demo")
	assert.Nil(t, err)
	reqHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/api":
			api.APIHandler(ctx)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}
	go fasthttp.ListenAndServe(":65002", reqHandler)
}

func TestMapBasicClientErrorOnWrongAddress(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:12345")
	client.SetUserAgent("JsonRPC Test Client")
	_, err := client.CallForMap("/api", "demo.Test", TestMapArgs{ID: "TESTER_ID_TestMapBasicClientErrorOnWrongAddress"})
	assert.Nil(t, err)
}

func TestMapBasicClientErrorOnAPIFormat(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://yandex.ru")
	client.SetUserAgent("JsonRPC Test Client")
	_, err := client.CallForMap("/api", "demo.Test", TestMapArgs{ID: "TESTER_ID_TestMapBasicClientErrorOnAPIFormat"})
	assert.Nil(t, err)
}

func TestMapBasicClientErrorOnAPIAnswer(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	_, err := client.CallForMap("/api", "demo.Test", TestMapArgs{ID: ""})
	assert.Nil(t, err)
}

func TestMapBasicClient(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuthHeader("user", "password")
	dstT, err := client.CallForMap("/api", "demo.Test", TestMapArgs{ID: "TESTER_ID_TestMapBasicClient"})
	assert.Nil(t, err)
	val, ok := dstT["log_id"]
	if assert.NotEqual(t, ok, false) {
		assert.Equal(t, "TESTER_ID_TestMapBasicClient", val)
	}
}

func TestMapBasicClientWithDefaultUserAgent(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstT, err := client.CallForMap("/api", "demo.Test", TestMapArgs{ID: "TESTER_ID_TestMapBasicClientWithDefaultUserAgent"})
	assert.Nil(t, err)
	val, ok := dstT["log_id"]
	if assert.NotEqual(t, ok, false) {
		assert.Equal(t, "TESTER_ID_TestMapBasicClientWithDefaultUserAgent", val)
	}
}

func TestMapLoggingDevNullClient(t *testing.T) {
	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
		Level:     logrus.DebugLevel,
	}
	client := NewClientWithLogger(logger)
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuthHeader("user", "password")
	dstT, err := client.CallForMap("/api", "demo.Test", TestMapArgs{ID: "TESTER_ID_TestMapLoggingDevNullClient"})
	assert.Nil(t, err)
	val, ok := dstT["log_id"]
	if assert.NotEqual(t, ok, false) {
		assert.Equal(t, "TESTER_ID_TestMapLoggingDevNullClient", val)
	}
}

func TestMapLoggingClient(t *testing.T) {
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
		Level:     logrus.DebugLevel,
	}
	client := NewClientWithLogger(logger)
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuthHeader("user", "password")
	dstT, err := client.CallForMap("/api", "demo.Test", TestMapArgs{ID: "TESTER_ID_TestMapLoggingClient"})
	assert.Nil(t, err)
	val, ok := dstT["log_id"]
	if assert.NotEqual(t, ok, false) {
		assert.Equal(t, "TESTER_ID_TestMapLoggingClient", val)
	}
}
