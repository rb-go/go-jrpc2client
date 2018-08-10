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
type DemoPointerAPI struct{}

type TestPointerArgs struct {
	ID string
}

type TestPointerReply struct {
	LogID     string `json:"log_id"`
	UserAgent string `json:"user_agent"`
}

// Test Method to test
func (h *DemoPointerAPI) Test(ctx *fasthttp.RequestCtx, args *TestPointerArgs, reply *TestPointerReply) error {
	if args.ID == "" {
		return &jrpc2server.Error{Code: jrpc2server.JErrorInvalidParams, Message: "ID should not be empty"}
	}
	reply.LogID = args.ID
	reply.UserAgent = string(ctx.Request.Header.UserAgent())
	return nil
}

// TestUserAgent Method to test user agent value
func (h *DemoPointerAPI) TestUserAgent(ctx *fasthttp.RequestCtx, args *TestPointerArgs, reply *TestPointerReply) error {
	reply.UserAgent = string(ctx.Request.Header.UserAgent())
	return nil
}

func TestPointerPrepare(t *testing.T) {
	api := jrpc2server.NewServer()
	err := api.RegisterService(new(DemoPointerAPI), "demo")
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

func TestPointerBasicClientErrorOnWrongAddress(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:12345")
	dstP := &TestPointerReply{}
	err := client.Call("/api", "demo.Test", TestPointerArgs{ID: ""}, dstP)
	assert.NotNil(t, err)
}

func TestPointerBasicClientErrorOnAPIFormat(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://yandex.ru")
	dstP := &TestPointerReply{}
	err := client.Call("/api", "demo.Test", TestPointerArgs{ID: ""}, dstP)
	assert.NotNil(t, err)
}

func TestPointerBasicClientErrorOnAPIAnwser(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestPointerReply{}
	err := client.Call("/api", "demo.Test", TestPointerArgs{ID: ""}, dstP)
	assert.NotNil(t, err)
}

func TestPointerDefaultUserAgentClient(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestPointerReply{}
	err := client.Call("/api", "demo.TestUserAgent", TestPointerArgs{ID: "TESTER_ID_TestPointerDefaultUserAgentClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, userAgent, dstP.UserAgent)
}

func TestPointerCustomUserAgentClient(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	dstP := &TestPointerReply{}
	err := client.Call("/api", "demo.TestUserAgent", TestPointerArgs{ID: "TESTER_ID_TestPointerCustomUserAgentClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "JsonRPC Test Client", dstP.UserAgent)
}

func TestPointerBasicAuthClient(t *testing.T) {
	client := NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuthHeader("user", "password")
	dstP := &TestPointerReply{}
	err := client.Call("/api", "demo.Test", TestPointerArgs{ID: "TESTER_ID_TestPointerBasicClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestPointerLoggingClient", dstP.LogID)
}

func TestPointerLoggingDevNullClient(t *testing.T) {
	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
		Level:     logrus.DebugLevel,
	}
	client := NewClientWithLogger(logger)
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestPointerReply{}
	err := client.Call("/api", "demo.Test", TestPointerArgs{ID: "TESTER_ID_TestPointerLoggingDevNullClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestPointerLoggingClient", dstP.LogID)
}

func TestPointerLoggingClient(t *testing.T) {
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
		Level:     logrus.DebugLevel,
	}
	client := NewClientWithLogger(logger)
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestPointerReply{}
	err := client.Call("/api", "demo.Test", TestPointerArgs{ID: "TESTER_ID_TestPointerLoggingClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestPointerLoggingClient", dstP.LogID)
}
