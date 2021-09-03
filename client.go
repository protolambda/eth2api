package eth2api

import (
	"context"
	"fmt"
)

type Client interface {
	Request(ctx context.Context, req PreparedRequest) Response
}

// simplified Query interface. No duplicate key entries
type Query map[string]interface{}

type PreparedRequest interface {
	// The type of request
	Method() ReqMethod

	// Body, optional. Returns nil if no body. Only used for PUT/POST methods.
	Body() interface{}

	// Path to request, including any path variable contents
	Path() string

	// Query params to add to the request, may be nil
	Query() Query
}

type Response interface {
	// Decode into destination type. May throw a decoding error.
	// Or throws DecodeNoContentErr if it was an error without returned value.
	// Call with nil to just close the response contents.
	// May only be called once.
	Decode(dest interface{}) (code uint, err error)

	// TODO: maybe expose headers?
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

func FmtGET(format string, data ...interface{}) PreparedRequest {
	return PlainGET(fmt.Sprintf(format, data...))
}

func QueryGET(query Query, path string) PreparedRequest {
	return &fullReq{method: GET, path: path, body: nil, query: query}
}

func FmtQueryGET(query Query, format string, data ...interface{}) PreparedRequest {
	return &fullReq{method: GET, path: fmt.Sprintf(format, data...), body: nil, query: query}
}

func BodyPOST(path string, body interface{}) PreparedRequest {
	return &fullReq{method: POST, path: path, body: body, query: nil}
}

func SimpleRequest(ctx context.Context, cli Client, req PreparedRequest, dest interface{}) (exists bool, err error) {
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(dest)
	exists = code != 404
	return
}

func MinimalRequest(ctx context.Context, cli Client, req PreparedRequest, dest interface{}) (err error) {
	resp := cli.Request(ctx, req)
	_, err = resp.Decode(dest)
	return
}

type ClientFunc func(ctx context.Context, req PreparedRequest) Response

func (fn ClientFunc) Request(ctx context.Context, req PreparedRequest) Response {
	return fn(ctx, req)
}

type ClientApiErr struct {
	error
}

func (ce ClientApiErr) Code() uint {
	return 400
}

type ClientErr struct {
	error
}

func (ce ClientErr) Decode(dest interface{}) (uint, error) {
	return 0, fmt.Errorf("client usage error, cannot decode: %w", ce.error)
}
