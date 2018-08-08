package jrpc2_client

import (
	"encoding/base64"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/erikdubbelboer/fasthttp"
)

func NewClient() *Client {
	return &Client{
		Logger: &logrus.Logger{
			Out:       os.Stdout,
			Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
			Level:     logrus.WarnLevel,
		},
	}
}

func NewClientWithLogger(logger *logrus.Logger) *Client {
	return &Client{
		Logger: logger,
	}
}

func (cl *Client) SetLogger(logger *logrus.Logger) {
	cl.Logger = logger
}

func (cl *Client) SetBaseURL(baseURL string) {
	cl.BaseUrl = baseURL
}

func (cl *Client) SetBasicAuth(login string, password string) {
	cl.Authentificate = "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password))
}

func (cl *Client) SetUserAgent(userAgent string) {
	cl.UserAgent = userAgent
}

func (cl *Client) Call(urlPath string, method string, args interface{}, dst interface{}) error {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(cl.BaseUrl + urlPath)

	setHeadersFromConfig(cl, req)

	req.Header.SetMethod("POST")
	byteBody, err := encodeClientRequest(method, args)
	if err != nil {
		return err
	}

	tmp := logrus.Fields{}
	tmp["headers"] = req.Header.String()
	tmp["request"] = byteBody
	debugLogging(cl, tmp, "request prepared")

	req.SetBody(byteBody)
	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	if err := client.Do(req, resp); err != nil {
		return err
	}

	tmp = logrus.Fields{}
	tmp["headers"] = req.Header.String()
	tmp["response"] = resp.Body()
	debugLogging(cl, tmp, "response received")

	return decodeClientResponse(resp.Body(), dst)
}

func (cl *Client) BatchCall(urlPath string, method string, args interface{}) {

}

func (cl *Client) AsyncCall(ch chan<- interface{}, urlPath string, method string, args interface{}) {
	var result interface{}
	ch <- result
}
