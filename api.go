package eth2api

import (
	"context"
	"errors"
	"fmt"
)

type ReqMethod uint

const (
	GET ReqMethod = iota
	PUT
	POST
)

// simplified Query interface. No duplicate key entries
type Query map[string]interface{}

type Request interface {
	// The type of request
	Method() ReqMethod

	// Body, optional. Returns nil if no body. Only used for PUT/POST methods.
	// Should only be retrieved once, to keep io with files, buffers etc. clean.
	Body() interface{}

	// Path to request, including any path variable contents
	Path() string

	// Query params to add to the request, may be nil
	Query() Query
}

var MissingRequiredParamErr = errors.New("missing required param")

var DecodeNoContentErr = errors.New("no contents were available to decode")

type Response interface {
	// Decode into destination type. May throw a decoding error.
	// Or throws DecodeNoContentErr if it was an error without returned value.
	// May only be called once.
	Decode(dest interface{}) error
	// Err when contents are not available / failed, always nil when code is 200.
	// For other 20x codes the error will be present, but may be ignored by the user.
	Err() ApiError
}

type Client interface {
	Request(ctx context.Context, req Request) Response
}

type fullReq struct {
	method ReqMethod
	path   string
	body   interface{}
	query  Query
}

func (p *fullReq) Method() ReqMethod {
	return p.method
}

func (p *fullReq) Body() interface{} {
	return p.body
}

func (p *fullReq) Path() string {
	return p.path
}

func (p *fullReq) Query() Query {
	return p.query
}

type PlainGET string

func (p PlainGET) Method() ReqMethod {
	return GET
}

func (p PlainGET) Body() interface{} {
	return nil
}

func (p PlainGET) Path() string {
	return string(p)
}

func (p PlainGET) Query() Query {
	return nil
}

func FmtGET(format string, data ...interface{}) Request {
	return PlainGET(fmt.Sprintf(format, data...))
}

func QueryGET(query Query, path string) Request {
	return &fullReq{method: POST, path: path, body: nil, query: query}
}

func FmtQueryGET(query Query, format string, data ...interface{}) Request {
	return &fullReq{method: POST, path: fmt.Sprintf(format, data...), body: nil, query: query}
}

func BodyPOST(path string, body interface{}) Request {
	return &fullReq{method: POST, path: path, body: body, query: nil}
}

func SimpleRequest(ctx context.Context, cli Client, req Request, dest interface{}) (exists bool, err error) {
	resp := cli.Request(ctx, req)
	if err := resp.Err(); err != nil {
		if err.Code() == 404 {
			return false, nil
		}
		return false, err
	}
	exists = true
	if dest != nil {
		err = resp.Decode(dest)
	}
	return
}

func MinimalRequest(ctx context.Context, cli Client, req Request, dest interface{}) (err error) {
	resp := cli.Request(ctx, req)
	if err := resp.Err(); err != nil {
		return err
	}
	if dest != nil {
		err = resp.Decode(dest)
	}
	return
}
