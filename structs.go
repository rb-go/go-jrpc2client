package jrpc2_client

import (
	"errors"

	"encoding/json"

	"github.com/Sirupsen/logrus"
)

type ErrorCode int

const (
	E_PARSE          ErrorCode = -32700
	E_INVALID_REQ    ErrorCode = -32600
	E_NO_METHOD      ErrorCode = -32601
	E_INVALID_PARAMS ErrorCode = -32602
	E_INTERNAL       ErrorCode = -32603
	E_SERVER         ErrorCode = -32000
)

var ErrNullResult = errors.New("result is null")

type clientConfig struct {
	UserAgent      string
	Authentificate string
	BaseUrl        string
	Logger         *logrus.Logger
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
	Id uint64 `json:"id"`
}

// clientResponse represents a JSON-RPC response returned to a client.
type clientResponse struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   *json.RawMessage `json:"error"`
}

type Error struct {
	// A Number that indicates the error type that occurred.
	Code ErrorCode `json:"code"` /* required */

	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"` /* required */

	// A Primitive or Structured value that contains additional information about the error.
	Data interface{} `json:"data"` /* optional */
}

func (e *Error) Error() string {
	return e.Message
}
