package jrpc2client

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/Sirupsen/logrus"
	"sync"
)

// ErrorCode type for error codes
type ErrorCode int

const (
	// JErrorParse Parse error - Invalid JSON was received by the server.
	// An error occurred on the server while parsing the JSON text.
	JErrorParse ErrorCode = -32700

	// JErrorInvalidReq Invalid Request - The JSON sent is not a valid Request object.
	JErrorInvalidReq ErrorCode = -32600

	// JErrorNoMethod Method not found - The method does not exist / is not available.
	JErrorNoMethod ErrorCode = -32601

	// JErrorInvalidParams Invalid params - Invalid method parameter(s).
	JErrorInvalidParams ErrorCode = -32602

	// JErrorInternal Internal error - Internal JSON-RPC error.
	JErrorInternal ErrorCode = -32603

	// JErrorServer Server error - Reserved for implementation-defined server-errors.
	JErrorServer ErrorCode = -32000

	userAgent string = "RIFTBIT-GOLANG-JRPC2-CLIENT"
)

// ErrNullResult it returns error when result answer is empty
var ErrNullResult = errors.New("result is null")

// Client basic struct that contains all method to work with JSON-RPC 2.0 protocol
type Client struct {
	userAgent         string
	BaseURL           string
	connectionTimeout time.Time
	writeTimeout      time.Time
	readTimeout       time.Time
	customHeaders     map[string]string
	clientPool        *sync.Pool
	logger            *logrus.Logger
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

// Error basic error struct to process API errors
type Error struct {
	// A Number that indicates the error type that occurred.
	Code ErrorCode `json:"code"` /* required */

	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"` /* required */

	// A Primitive or Structured value that contains additional information about the error.
	Data interface{} `json:"data"` /* optional */
}

// Error returns string based error
func (e *Error) Error() string {
	return e.Message
}
