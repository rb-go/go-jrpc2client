package jrpc2client

import (
	"math/rand"

	"reflect"

	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/pquerna/ffjson/ffjson"
)

func printObject(v interface{}) string {
	res2B, _ := ffjson.Marshal(v)
	return string(res2B)
}

func debugLogging(clientCfg *Client, fields logrus.Fields, message string) {
	if clientCfg.logger.Level == logrus.DebugLevel {
		for i, v := range fields {
			if reflect.TypeOf(v).String() == "[]uint8" {
				fields[i] = strings.Split(string(v.([]uint8)), "\r\n")
			}

			if i == "headers" && reflect.TypeOf(v).String() == "string" {
				fields[i] = strings.Split(fields[i].(string), "\r\n")
			}
		}
		clientCfg.logger.WithFields(fields).Debugln(message)
	}
}

// encodeClientRequest encodes parameters for a JSON-RPC client request.
func encodeClientRequest(method string, args interface{}) ([]byte, error) {
	c := &clientRequest{
		ID:      uint64(rand.Int63()),
		Version: "2.0",
		Method:  method,
		Params:  args,
	}
	return ffjson.Marshal(c)
}

// decodeClientResponse decodes the response body of a client request into the interface reply.
func decodeClientResponse(r []byte) (interface{}, error) {
	var c clientResponse
	if err := ffjson.NewDecoder().Decode(r, &c); err != nil {
		return nil, &Error{Code: JErrorParse, Message: err.Error()}
	}
	if c.Error != nil {
		jsonErr := &Error{}
		if err := ffjson.Unmarshal(*c.Error, jsonErr); err != nil {
			return nil, &Error{Code: JErrorInternal, Message: string(*c.Error)}
		}
		return nil, jsonErr
	}
	if c.Result == nil {
		return nil, &Error{Code: JErrorServer, Message: ErrNullResult.Error()}
	}
	var dst interface{}
	if err := ffjson.Unmarshal(*c.Result, &dst); err != nil {
		return nil, &Error{Code: JErrorInternal, Message: ErrNullResult.Error()}
	}
	return dst, nil
}
