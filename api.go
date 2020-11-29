package eth2api

import (
	"context"
	"errors"
)

type Request interface {
	// Path to request, including any path variable contents
	Path() string

	// Query params to add to the request, may be nil
	Query() map[string]interface{}
}

var DecodeNoContentErr = errors.New("no contents were available to decode")

type Response interface {
	// Decode into destination type. May throw a decoding error.
	// Or throws DecodeNoContentErr if it was an error without returned value.
	// May only be called once.
	Decode(dest interface{}) error
	// when contents are not available, i.e. not 200
	Err() *ErrorMessage
}

type Client interface {
	Request(ctx context.Context, req Request) Response
}

type PlainRequest string

func (p PlainRequest) Path() string {
	return string(p)
}

func (p PlainRequest) Query() map[string]interface{} {
	return nil
}
