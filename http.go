package eth2api

import (
	"bytes"
	"context"
	"encoding"
	"fmt"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Eth2HttpClient struct {
	Addr  string
	Cli   HTTPClient
	Codec Codec
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
		return &HttpResponse{Response: resp, Codec: cli.Codec}
	case POST:
		var buf bytes.Buffer
		if err := cli.Codec.EncodeRequestBody(&buf, req.Body()); err != nil {
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
		return &HttpResponse{Response: resp, Codec: cli.Codec}
	default:
		return ClientErr{fmt.Errorf("unrecognized request method enum value: %s", method)}
	}
}

type HttpResponse struct {
	*http.Response
	Codec Codec
}

func (resp *HttpResponse) Decode(dest interface{}) (code uint, err error) {
	hr := resp.Response
	code = uint(hr.StatusCode)
	err = resp.Codec.DecodeResponseBody(code, hr.Body, dest)
	return
}

type HttpRouter struct {
	httprouter.Router
	Codec         Codec
	OnEncodingErr func(error)
}

func NewHttpRouter() *HttpRouter {
	// TODO: add panic and notfound handlers
	return &HttpRouter{
		Router: httprouter.Router{
			RedirectTrailingSlash:  true,
			RedirectFixedPath:      true,
			HandleMethodNotAllowed: true,
			HandleOPTIONS:          true,
		},
		Codec: JSONCodec{},
	}
}

var _ http.Handler = (*HttpRouter)(nil)

func (r *HttpRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

type httpRequest struct {
	req    *http.Request
	query  url.Values
	params httprouter.Params
	codec  Codec
}

func (req httpRequest) DecodeBody(dst interface{}) error {
	return req.codec.DecodeRequestBody(req.req.Body, dst)
}

func (req httpRequest) Param(name string) string {
	return req.params.ByName(name)
}

func (req httpRequest) Query(name string) (values []string, ok bool) {
	values, ok = req.query[name]
	return
}

func (r *HttpRouter) AddRoute(route Route) {
	r.Router.Handle(string(route.Method()), route.Route(),
		func(respw http.ResponseWriter, req *http.Request, params httprouter.Params) {
			resp := route.Handle(req.Context(), httpRequest{
				req:    req,
				query:  req.URL.Query(),
				params: params,
				codec:  r.Codec,
			})
			h := respw.Header()
			for k, v := range resp.Headers() {
				h.Add(k, v)
			}
			respw.WriteHeader(int(resp.Code()))
			if err := r.Codec.EncodeResponseBody(respw, resp.Body()); err != nil && r.OnEncodingErr != nil {
				r.OnEncodingErr(err)
			}
		},
	)
}
