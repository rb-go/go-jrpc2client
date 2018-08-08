package jrpc2_client

import (
	"math/rand"

	"reflect"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/erikdubbelboer/fasthttp"
	"github.com/pquerna/ffjson/ffjson"
)

func debugLogging(clientCfg *ClientConfig, fields logrus.Fields, message string) {
	if clientCfg.Logger.Level == logrus.DebugLevel {

		for i, v := range fields {
			if reflect.TypeOf(v).String() == "[]uint8" {
				fields[i] = strings.Split(string(v.([]uint8)), "\r\n")
			}

			if i == "headers" && reflect.TypeOf(v).String() == "string" {
				fields[i] = strings.Split(fields[i].(string), "\r\n")
			}

		}
		clientCfg.Logger.WithFields(fields).Debugln(message)
	}
}

func setHeadersFromConfig(clientCfg *ClientConfig, req *fasthttp.Request) {
	if clientCfg.UserAgent != "" {
		req.Header.SetUserAgent(clientCfg.UserAgent)
	} else {
		req.Header.SetUserAgent(clientCfg.UserAgent)
	}

	if clientCfg.Authentificate != "" {
		req.Header.Set("Authorization", clientCfg.Authentificate)
	}
}

// encodeClientRequest encodes parameters for a JSON-RPC client request.
func encodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &clientRequest{
		Id:      uint64(rand.Int63()),
		Version: "2.0",
		Method:  method,
		Params:  args,
	}
	return ffjson.Marshal(c)
}

// decodeClientResponse decodes the response body of a client request into
// the interface reply.
func decodeClientResponse(r []byte, dst interface{}) error {
	var c clientResponse
	if err := ffjson.NewDecoder().Decode(r, &c); err != nil {
		return &Error{
			Code:    E_PARSE,
			Message: err.Error(),
		}
	}
	if c.Error != nil {
		jsonErr := &Error{}
		if err := ffjson.Unmarshal(*c.Error, jsonErr); err != nil {
			return &Error{
				Code:    E_INTERNAL,
				Message: string(*c.Error),
			}
		}
		return jsonErr
	}
	if c.Result == nil {
		return &Error{
			Code:    E_SERVER,
			Message: ErrNullResult.Error(),
		}
	}
	if err := ffjson.Unmarshal(*c.Result, &dst); err != nil {
		return &Error{
			Code:    E_INTERNAL,
			Message: ErrNullResult.Error(),
		}
	}
	return nil
}
