package jrpc2client

import (
	"encoding/base64"
	"os"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/erikdubbelboer/fasthttp"
)

func getDefaultHeadersMap() map[string]string {
	headers := make(map[string]string)
	headers["User-Agent"] = userAgent
	return headers
}

func createNewClient(logger *logrus.Logger) *Client {
	return &Client{
		clientPool: &sync.Pool{
			New: func() interface{} {
				return new(fasthttp.Client)
			},
		},
		customHeaders: getDefaultHeadersMap(),
		logger:        logger,
	}
}

// NewClient returns new configured Client to start work with JSON-RPC 2.0 protocol
func NewClient() *Client {
	return createNewClient(&logrus.Logger{Out: os.Stdout, Formatter: &logrus.JSONFormatter{DisableTimestamp: false}, Level: logrus.WarnLevel})
}

// NewClientWithLogger returns new configured Client with custom Logger configureation (based on Sirupsen/logrus) to start work with JSON-RPC 2.0 protocol
func NewClientWithLogger(logger *logrus.Logger) *Client {
	return createNewClient(logger)
}

// SetBaseURL setting basic url for API
func (cl *Client) SetBaseURL(baseURL string) {
	cl.BaseURL = baseURL
}

// SetClientTimeout this method sets globally for client its timeout
func (cl *Client) SetClientTimeout(duration time.Duration) {
	cl.clientTimeout = duration
}

// SetCustomHeader setting custom header
func (cl *Client) SetCustomHeader(headerName string, headerValue string) {
	cl.customHeaders[headerName] = headerValue
}

// DeleteCustomHeader delete custom header
func (cl *Client) DeleteCustomHeader(headerName string) {
	delete(cl.customHeaders, headerName)
}

// SetBasicAuthHeader setting basic auth header
func (cl *Client) SetBasicAuthHeader(login string, password string) {
	cl.SetCustomAuthHeader("Basic", base64.StdEncoding.EncodeToString([]byte(login+":"+password)))
}

// SetCustomAuthHeader setting custom auth header with type of auth and auth data
func (cl *Client) SetCustomAuthHeader(authType string, authData string) {
	cl.SetCustomHeader("Authorization", authType+" "+authData)
}

// DeleteAuthHeader clear basic auth header
func (cl *Client) DeleteAuthHeader() {
	cl.DeleteCustomHeader("Authorization")
}

// SetUserAgent setting custom User Agent header
func (cl *Client) SetUserAgent(userAgent string) {
	cl.SetCustomHeader("User-Agent", userAgent)
}

func (cl *Client) makeCallRequest(urlPath string, method string, args interface{}) ([]byte, error) {
	req := fasthttp.AcquireRequest()
	defer req.Reset()
	req.SetRequestURI(cl.BaseURL + urlPath)

	for key, val := range cl.customHeaders {
		req.Header.Set(key, val)
	}

	req.Header.SetMethod("POST")
	byteBody, err := encodeClientRequest(method, args)
	if err != nil {
		return nil, err
	}

	debugLogging(cl, logrus.Fields{"headers": req.Header.String(), "request": byteBody}, "request prepared")

	req.SetBody(byteBody)
	resp := fasthttp.AcquireResponse()
	defer resp.Reset()

	client := cl.clientPool.Get().(*fasthttp.Client)

	if cl.clientTimeout == 0 {
		if err := client.Do(req, resp); err != nil {
			return nil, err
		}
	} else {
		if err := client.DoTimeout(req, resp, cl.clientTimeout); err != nil {
			return nil, err
		}
	}

	cl.clientPool.Put(client)
	debugLogging(cl, logrus.Fields{"headers": req.Header.String(), "response": resp.Body()}, "response received")
	return resp.Body(), nil
}

// Call run remote procedure on JSON-RPC 2.0 API with parsing answer to provided structure or interface
func (cl *Client) Call(urlPath string, method string, args interface{}, dst interface{}) error {
	resp, err := cl.makeCallRequest(urlPath, method, args)
	if err != nil {
		return err
	}
	err = decodeClientResponse(resp, &dst)
	return err
}

// CallForMap run remote procedure on JSON-RPC 2.0 API with returning map[string]interface{}
func (cl *Client) CallForMap(urlPath string, method string, args interface{}) (map[string]interface{}, error) {
	resp, err := cl.makeCallRequest(urlPath, method, args)
	if err != nil {
		return nil, err
	}
	dst := make(map[string]interface{})
	err = decodeClientResponse(resp, &dst)
	return dst, err
}

/*
func (cl *Client) CallBatch(urlPath string, method string, args interface{}) {

}

func (cl *Client) AsyncCall(urlPath string, method string, args interface{}, ch chan<- interface{}) {
	var result interface{}
	ch <- result
}
*/
