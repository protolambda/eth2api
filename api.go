package eth2api

import (
	"io"
)

type Codec interface {
	EncodeResponseBody(w io.Writer, data interface{}) error
	DecodeResponseBody(code uint, r io.ReadCloser, dest interface{}) error
	EncodeRequestBody(w io.Writer, body interface{}) error
	DecodeRequestBody(r io.ReadCloser, dst interface{}) error
	ContentType() []string
}

type ReqMethod string

const (
	GET  ReqMethod = "GET"
	POST ReqMethod = "POST"
)

// DataWrap is a util to accommodate responses which are wrapped
// with a single field container with key "data".
type DataWrap struct {
	Data interface{} `json:"data"`
}

func Wrap(data interface{}) *DataWrap {
	return &DataWrap{Data: data}
}
