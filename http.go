package eth2api

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Eth2HttpClient struct {
	Addr string
	Cli  HTTPClient
}

func (cli *Eth2HttpClient) Request(ctx context.Context, req PreparedRequest) Response {
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

type HttpResponse http.Response

func (resp *HttpResponse) Decode(dest interface{}) (code uint, err error) {
	hr := (*http.Response)(resp)
	return DecodeBody(uint(hr.StatusCode), hr.Body, dest)
}

type HttpRouter struct {
}

var _ http.Handler = (*HttpRouter)(nil)

func (r *HttpRouter) ServeHTTP(http.ResponseWriter, *http.Request) {
	// TODO: use the constructed router
}

func (r *HttpRouter) AddRoute(route Route) {
	// TODO: use julienschmidt/httprouter to mux & extract vars from route path
	//handle := http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
	//	route.Handle()
	//})
}
