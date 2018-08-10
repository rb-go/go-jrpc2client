package jrpc2client_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/riftbit/jrpc2client"
	"github.com/riftbit/jrpc2server"
	"github.com/stretchr/testify/assert"
)

// DemoAPI area
type DemoAPI struct{}

type TestArgs struct {
	ID string
}

type TestReply struct {
	LogID     string `json:"log_id"`
	UserAgent string `json:"user_agent"`
}

// Test Method to test
func (h *DemoAPI) Test(ctx *fasthttp.RequestCtx, args *TestArgs, reply *TestReply) error {
	if args.ID == "" {
		return &jrpc2server.Error{Code: jrpc2server.JErrorInvalidParams, Message: "ID should not be empty"}
	}
	reply.LogID = args.ID
	reply.UserAgent = string(ctx.Request.Header.UserAgent())
	return nil
}

// TestUserAgent Method to test user agent value
func (h *DemoAPI) TestUserAgent(ctx *fasthttp.RequestCtx, args *TestArgs, reply *TestReply) error {
	reply.UserAgent = string(ctx.Request.Header.UserAgent())
	return nil
}

// TestUserAgent Method to test user agent value
func (h *DemoAPI) TestClientTimeout(ctx *fasthttp.RequestCtx, args *TestArgs, reply *TestReply) error {
	time.Sleep(100 * time.Millisecond)
	reply.UserAgent = string(ctx.Request.Header.UserAgent())
	return nil
}

func init() {
	api := jrpc2server.NewServer()
	err := api.RegisterService(new(DemoAPI), "demo")
	if err != nil {
		panic(err)
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

func TestClient_Call_WrongAddress(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:12345")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.Test", TestArgs{ID: ""}, dstP)
	assert.NotNil(t, err)
}

func TestClient_Call_WrongAnswerFormat(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://yandex.ru")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.Test", TestArgs{ID: ""}, dstP)
	assert.NotNil(t, err)
}

func TestClient_Call_WrongAPIAnswer(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.Test", TestArgs{ID: ""}, dstP)
	assert.NotNil(t, err)
}

func TestClient_CallForMap(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuthHeader("user", "password")
	dstT, err := client.CallForMap("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestMapBasicClient"})
	assert.Nil(t, err)
	val, ok := dstT["log_id"]
	if assert.NotEqual(t, ok, false) {
		assert.Equal(t, "TESTER_ID_TestMapBasicClient", val)
	}
}

func TestClient_CallForMap_WrongAddress(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:12345")
	client.SetUserAgent("JsonRPC Test Client")
	client.SetBasicAuthHeader("user", "password")
	_, err := client.CallForMap("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestMapBasicClient"})
	assert.NotNil(t, err)
}

func TestClient_SetUserAgent_default(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.TestUserAgent", TestArgs{ID: "TESTER_ID_TestDefaultUserAgentClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "RIFTBIT-JRPC2-CLIENT", dstP.UserAgent)
}

func TestClient_SetUserAgent_custom(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetUserAgent("JsonRPC Test Client")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.TestUserAgent", TestArgs{ID: "TESTER_ID_TestCustomUserAgentClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "JsonRPC Test Client", dstP.UserAgent)
}

func TestClient_SetBasicAuthHeader(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	client.SetBasicAuthHeader("user", "password")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestBasicAuthClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestBasicAuthClient", dstP.LogID)
}

func TestNewClientWithLogger_devNull(t *testing.T) {
	logger := &logrus.Logger{
		Out:       ioutil.Discard,
		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
		Level:     logrus.DebugLevel,
	}
	client := jrpc2client.NewClientWithLogger(logger)
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestLoggingDevNullClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestLoggingDevNullClient", dstP.LogID)
}

func TestNewClientWithLogger(t *testing.T) {
	logger := &logrus.Logger{
		Out:       os.Stdout,
		Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
		Level:     logrus.DebugLevel,
	}
	client := jrpc2client.NewClientWithLogger(logger)
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestLoggingClient"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestLoggingClient", dstP.LogID)
}

func TestClient_Call_double(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}

	err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestDoubleCallBasicClient_1"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestDoubleCallBasicClient_1", dstP.LogID)

	err = client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestDoubleCallBasicClient_2"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestDoubleCallBasicClient_2", dstP.LogID)

}

func TestClient_Call_tripple(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}

	err := client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestTrippleCallBasicClient_1"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestTrippleCallBasicClient_1", dstP.LogID)

	err = client.Call("/api", "demo.Test", TestArgs{ID: "TESTER_ID_TestTrippleCallBasicClient_2"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "TESTER_ID_TestTrippleCallBasicClient_2", dstP.LogID)

	client.SetUserAgent("JsonRPC Test Client")
	err = client.Call("/api", "demo.TestUserAgent", TestArgs{ID: "TESTER_ID_TestTrippleCallBasicClient_3"}, dstP)
	assert.Nil(t, err)
	assert.Equal(t, "JsonRPC Test Client", dstP.UserAgent)
}

func TestClient_SetClientTimeout_expectedErr(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}
	client.SetClientTimeout(10 * time.Millisecond)
	err := client.Call("/api", "demo.TestClientTimeout", TestArgs{ID: ""}, dstP)
	assert.NotNil(t, err)
}

func TestClient_SetClientTimeout_success(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}
	client.SetClientTimeout(1 * time.Second)
	err := client.Call("/api", "demo.TestClientTimeout", TestArgs{ID: ""}, dstP)
	assert.Nil(t, err)
}

func TestClient_DeleteCustomHeaderClient(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}
	client.DeleteCustomHeader("User-Agent")
	err := client.Call("/api", "demo.TestClientTimeout", TestArgs{ID: ""}, dstP)
	assert.Nil(t, err)
}

func TestError_Error(t *testing.T) {
	client := jrpc2client.NewClient()
	client.SetBaseURL("http://127.0.0.1:65001")
	dstP := &TestReply{}
	err := client.Call("/api", "demo.Test", TestArgs{ID: ""}, dstP)
	assert.Equal(t, "ID should not be empty", err.Error())
}
