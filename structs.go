package jrpc2client

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// ErrorCode type for error codes
type ErrorCode int

const (
	userAgent          = "RIFTBIT-JRPC2-CLIENT"
	defaultContentType = "application/json"
)

// Client basic struct that contains all method to work with JSON-RPC 2.0 protocol
type Client struct {
	BaseURL       string
	clientTimeout time.Duration
	customHeaders map[string]string
	clientPool    *sync.Pool
	logger        *logrus.Logger
}

// clientRequest represents a JSON-RPC request sent by a client.
type clientRequest struct {
	// JSON-RPC protocol.
	Version string `json:"jsonrpc"`

	// A String containing the name of the method to be invoked.
	Method string `json:"method"`

	// Object to pass as request parameter to the method.
	Params interface{} `json:"params"`

	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	ID uint64 `json:"id"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type clientResponse struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   *json.RawMessage `json:"error"`
}
