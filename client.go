package jrpc2_client

import (
	"encoding/base64"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/erikdubbelboer/fasthttp"
)

func NewClient() *ClientConfig {
	return &ClientConfig{
		Logger: &logrus.Logger{
			Out:       os.Stdout,
			Formatter: &logrus.JSONFormatter{DisableTimestamp: false},
			Level:     logrus.WarnLevel,
		},
	}
}

func NewClientWithLogger(logger *logrus.Logger) *ClientConfig {
	return &ClientConfig{
		Logger: logger,
	}
}

func (cl *ClientConfig) SetLogger(logger *logrus.Logger) {
	cl.Logger = logger
}

func (cl *ClientConfig) SetBaseURL(baseURL string) {
	cl.BaseUrl = baseURL
}

func (cl *ClientConfig) SetBasicAuth(login string, password string) {
	cl.Authentificate = "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password))
}

func (cl *ClientConfig) SetUserAgent(userAgent string) {
	cl.UserAgent = userAgent
}

func (cl *ClientConfig) Call(urlPath string, method string, args interface{}, dst interface{}) error {
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

func (cl *ClientConfig) BatchCall(urlPath string, method string, args interface{}) {

}

func (cl *ClientConfig) AsyncCall(ch chan<- interface{}, urlPath string, method string, args interface{}) {
	var result interface{}
	ch <- result
}