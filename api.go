package eth2api

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ReqMethod uint

const (
	GET ReqMethod = iota
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

// DataWrap is a util to accommodate responses which are wrapped
// with a single field container with key "data".
type DataWrap struct {
	Data interface{} `json:"data"`
}

func Wrap(data interface{}) *DataWrap {
	return &DataWrap{Data: data}
}

type Response interface {
	// Decode into destination type. May throw a decoding error.
	// Or throws DecodeNoContentErr if it was an error without returned value.
	// Call with nil to just close the response contents.
	// May only be called once.
	Decode(dest interface{}) (code uint, err error)
}

type ClientFunc func(ctx context.Context, req Request) Response

func (fn ClientFunc) Request(ctx context.Context, req Request) Response {
	return fn(ctx, req)
}

type Client interface {
	Request(ctx context.Context, req Request) Response
}

type HttpResponse http.Response

func DecodeBody(code uint, body io.ReadCloser, dest interface{}) (codeOut uint, err error) {
	defer body.Close()
	codeOut = code
	if code < 200 {
		return code, fmt.Errorf("unexpected response status code: %d", code)
	} else if code < 300 {
		dec := json.NewDecoder(body)
		return code, dec.Decode(dest)
	} else {
		// TODO: handle indexed errors

		var ierr ErrorResponse
		dec := json.NewDecoder(body)
		if err := dec.Decode(&ierr); err != nil {
			return code, ClientApiErr{fmt.Errorf("failed to decode error response with status code: %d", code)}
		}
		if code < 500 {
			return code, &InvalidRequest{ierr}
		}
		if code == 503 {
			return code, &CurrentlySyncing{ierr}
		}
		if code < 600 {
			return code, &InternalError{ierr}
		}
		return code, &ierr
	}
	// TODO: could support more than just JSON by looking at Content-Type,
	// and using Content-Length for fast SSZ streaming
	// (after unwrapping the contents from the inner Data field and checking SSZ support,
	//  and sourcing a spec from somewhere)
}

func (resp *HttpResponse) Decode(dest interface{}) (code uint, err error) {
	hr := (*http.Response)(resp)
	return DecodeBody(uint(hr.StatusCode), hr.Body, dest)
}

type HttpClient struct {
	Addr string
	Cli  *http.Client
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

func (cli *HttpClient) Request(ctx context.Context, req Request) Response {
	path := cli.Addr + req.Path()
	if q := req.Query(); q != nil {
		b := make(url.Values)
		for k, v := range req.Query() {
			if s, ok := v.(string); ok {
				b.Set(k, s)
			} else if sv, ok := v.(fmt.Stringer); ok {
				b.Set(k, sv.String())
			} else if tm, ok := v.(encoding.TextMarshaler); ok {
				tb, err := tm.MarshalText()
				if err != nil {
					return ClientErr{fmt.Errorf("failed to encode query key %s: %w", k, err)}
				}
				b.Set(k, string(tb))
			} else {
				return ClientErr{fmt.Errorf("failed to encode query key '%s': unknown type", k)}
			}
		}
		path += "?" + b.Encode()
	}
	method := req.Method()
	switch method {
	case GET:
		req, err := http.NewRequestWithContext(ctx, "GET", path, nil)
		if err != nil {
			return ClientErr{fmt.Errorf("failed to build GET request: %w", err)}
		}
		resp, err := cli.Cli.Do(req)
		if err != nil {
			return ClientErr{fmt.Errorf("failed to execute GET request: %w", err)}
		}
		return (*HttpResponse)(resp)
	case POST:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf) // TODO: different content-types
		if err := enc.Encode(req.Body()); err != nil {
			return ClientErr{fmt.Errorf("failed to encode POST request body: %w", err)}
		}
		req, err := http.NewRequestWithContext(ctx, "POST", path, &buf)
		if err != nil {
			return ClientErr{fmt.Errorf("failed to build POST request: %w", err)}
		}
		resp, err := cli.Cli.Do(req)
		if err != nil {
			return ClientErr{fmt.Errorf("failed to execute POST request: %w", err)}
		}
		return (*HttpResponse)(resp)
	default:
		return ClientErr{fmt.Errorf("unrecognized request method enum value: %d", method)}
	}
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
	var code uint
	code, err = resp.Decode(dest)
	exists = code != 404
	return
}

func MinimalRequest(ctx context.Context, cli Client, req Request, dest interface{}) (err error) {
	resp := cli.Request(ctx, req)
	_, err = resp.Decode(dest)
	return
}
