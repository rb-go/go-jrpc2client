package jrpc2client

import (
	"encoding/base64"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/erikdubbelboer/fasthttp"
)

// NewClient returns new configured Client to start work with JSON-RPC 2.0 protocol
func NewClient() *Client {
	return &Client{
		customHeaders: make(map[string]string),
		logger: &logrus.Logger{
			Out:       os.Stdout,
			Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
			Level:     logrus.WarnLevel,
		},
	}
}

// NewClientWithLogger returns new configured Client with custom Logger configureation (based on Sirupsen/logrus) to start work with JSON-RPC 2.0 protocol
func NewClientWithLogger(logger *logrus.Logger) *Client {
	return &Client{
		logger: logger,
	}
}

// SetBaseURL setting basic url for API
func (cl *Client) SetBaseURL(baseURL string) {
	cl.BaseURL = baseURL
}

// SetBasicAuth setting basic auth header
func (cl *Client) SetBasicAuth(login string, password string) {
	cl.customHeaders["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password))
}

// SetUserAgent setting custom User Agent header
func (cl *Client) SetUserAgent(userAgent string) {
	cl.customHeaders["User-Agent"] = userAgent
}

// Call run remote procedure on JSON-RPC 2.0 API
func (cl *Client) Call(urlPath string, method string, args interface{}) (interface{}, error) {
	req := fasthttp.AcquireRequest()
	defer req.Reset()

	req.SetRequestURI(cl.BaseURL + urlPath)
	req.Header.SetUserAgent(userAgent)

	for key, val := range cl.customHeaders {
		req.Header.Set(key, val)
	}

	req.Header.SetMethod("POST")
	byteBody, err := encodeClientRequest(method, args)
	if err != nil {
		return nil, err
	}

	tmp := logrus.Fields{}
	tmp["headers"] = req.Header.String()
	tmp["request"] = byteBody
	debugLogging(cl, tmp, "request prepared")

	req.SetBody(byteBody)
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return nil, err
	}

	tmp = logrus.Fields{}
	tmp["headers"] = req.Header.String()
	tmp["response"] = resp.Body()
	debugLogging(cl, tmp, "response received")

	return decodeClientResponse(resp.Body())
}

/*
func (cl *Client) BatchCall(urlPath string, method string, args interface{}) {

}

func (cl *Client) AsyncCall(urlPath string, method string, args interface{}, ch chan<- interface{}) {
	var result interface{}
	ch <- result
}
*/
