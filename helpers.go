package jrpc2client

import (
	"math/rand"
	"reflect"
	"strings"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/riftbit/jrpc2errors"
	"github.com/sirupsen/logrus"
)

// func printObject(v interface{}) string {
// 	res2B, _ := ffjson.Marshal(v)
// 	return string(res2B)
// }

// debugLogging method to show debug message with pre-processd values
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
		ID:      rand.Uint64(),
		Version: "2.0",
		Method:  method,
		Params:  args,
	}
	return ffjson.Marshal(c)
}

// decodeClientResponse decodes the response body of a client request into the interface reply.
func decodeClientResponse(r []byte, dst interface{}) error {
	var c clientResponse
	if err := ffjson.NewDecoder().Decode(r, &c); err != nil {
		return &jrpc2errors.Error{Code: jrpc2errors.ParseError, Message: err.Error()}
	}
	if c.Error != nil {
		jsonErr := &jrpc2errors.Error{}
		if err := ffjson.Unmarshal(*c.Error, jsonErr); err != nil {
			return &jrpc2errors.Error{Code: jrpc2errors.InternalError, Message: string(*c.Error)}
		}
		return jsonErr
	}
	if c.Result == nil {
		return &jrpc2errors.Error{Code: jrpc2errors.ServerError, Message: jrpc2errors.ErrNullResult.Error()}
	}

	if err := ffjson.Unmarshal(*c.Result, &dst); err != nil {
		return &jrpc2errors.Error{Code: jrpc2errors.InternalError, Message: jrpc2errors.ErrNullResult.Error()}
	}
	return nil
}
