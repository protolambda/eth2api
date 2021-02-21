package eth2api

import (
	"errors"
)

type ReqMethod uint

const (
	GET ReqMethod = iota
	POST
)

var MissingRequiredParamErr = errors.New("missing required param")

// DataWrap is a util to accommodate responses which are wrapped
// with a single field container with key "data".
type DataWrap struct {
	Data interface{} `json:"data"`
}

func Wrap(data interface{}) *DataWrap {
	return &DataWrap{Data: data}
}
