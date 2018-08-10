package jrpc2client

import (
	"io/ioutil"
	"os"
	"testing"

	"errors"

	"github.com/Sirupsen/logrus"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/riftbit/jrpc2server"
)

// DemoAPI area
type DemoAPI struct{}

type TestArgs struct {
	ID string
}

type TestReply struct {
	LogID string `json:"log_id"`
}

// Test Method to test
func (h *DemoAPI) Test(ctx *fasthttp.RequestCtx, args *TestArgs, reply *TestReply) error {
	if args.ID == "" {
		return &jrpc2server.Error{Code: jrpc2server.JErrorInvalidParams, Message: "ID should not be empty"}
	}
	reply.LogID = args.ID
	return nil
}

func TestPrepare(t *testing.T) {

	api := jrpc2server.NewServer()
	err := api.RegisterService(new(DemoAPI), "demo")

	if err != nil {
		t.Error(err)
	}

	reqHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/api":
			api.APIHandler(ctx)
		default:
			ctx.Error("Unsupported path", fasthttp.StatusNotFound)
		}
	}

	go fasthttp.ListenAndServe(":65001", reqHandler)
}

func TestBasicClientErrorOnWrongAddress(t *testing.T) {

	client := NewClient()

	client.SetBaseURL("http://127.0.0.1:12345")
	client.SetUserAgent("JsonRPC Test Client")

	_, err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestBasicClientErrorOnWrongAddress"})
	if err == nil {
		t.Error(errors.New("expected error but not received"))
	}
}

func TestBasicClientErrorOnAPIFormat(t *testing.T) {

	client := NewClient()

	client.SetBaseURL("http://yandex.ru")
	client.SetUserAgent("JsonRPC Test Client")

	_, err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestBasicClientErrorOnAPIFormat"})
	if err == nil {
		t.Error(errors.New("expected error but not received"))
	}
}

func TestBasicClientErrorOnAPIAnser(t *testing.T) {

	client := NewClient()

	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")

	_, err := client.Call("/api", "demo.Test", TestArgs{ID: ""})
	if err == nil {
		t.Error(errors.New("expected error but not received"))
	}
}

func TestBasicClient(t *testing.T) {

	client := NewClient()

	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuth("user", "password")

	dstT, err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestBasicClient"})
	if err != nil {
		t.Error(err)
	}
	if dstT.(TestReply).LogID != "TESTER_ID_TestBasicClient" {
		t.Error("unexpected answer in LogID")
	}
}

func TestBasicClientWithDefaultUserAgent(t *testing.T) {

	client := NewClient()

	client.SetBaseURL("http://127.0.0.1:65001")

	dstT, err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestBasicClientWithDefaultUserAgent"})
	if err != nil {
		t.Error(err)
	}
	if dstT.(TestReply).LogID != "TESTER_ID_TestBasicClientWithDefaultUserAgent" {
		t.Error("unexpected answer in LogID")
	}
}

func TestLoggingDevNullClient(t *testing.T) {
	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
		Level:     logrus.DebugLevel,
	}

	client := NewClientWithLogger(logger)

	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuth("user", "password")

	dstT, err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestLoggingDevNullClient"})
	if err != nil {
		t.Error(err)
	}
	if dstT.(TestReply).LogID != "TESTER_ID_TestLoggingDevNullClient" {
		t.Error("unexpected answer in LogID")
	}
}

func TestLoggingClient(t *testing.T) {
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
		Level:     logrus.DebugLevel,
	}

	client := NewClientWithLogger(logger)

	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuth("user", "password")

	dstT, err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestLoggingClient"})
	if err != nil {
		t.Error(err)
	}
	if dstT.(TestReply).LogID != "TESTER_ID_TestLoggingClient" {
		t.Error("unexpected answer in LogID")
	}
}
