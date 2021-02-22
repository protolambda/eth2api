package eth2api

import "context"

type Server interface {
	AddRoute(handler Route)
}

type Route interface {
	Method() ReqMethod
	Route() string
	Handle(ctx context.Context, req Request) PreparedResponse
}

type Request interface {
	DecodeBody(dst interface{}) error
	Param(name string) string
	Query(name string) (values []string, ok bool)

	// TODO: maybe expose headers?
}

type HandlerFn func(ctx context.Context, req Request) PreparedResponse

type route struct {
	method  ReqMethod
	pattern string
	handle  HandlerFn
}

func (r *route) Method() ReqMethod {
	return r.method
}

func (r *route) Route() string {
	return r.pattern
}

func (r *route) Handle(ctx context.Context, req Request) PreparedResponse {
	return r.handle(ctx, req)
}

func MakeRoute(method ReqMethod, path string, handle HandlerFn) Route {
	return &route{method, path, handle}
}
