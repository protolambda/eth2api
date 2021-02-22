package eth2api

import (
	"encoding/json"
	"fmt"
	"io"
)

type JSONCodec struct{}

func (JSONCodec) EncodeResponseBody(w io.Writer, data interface{}) error {
	enc := json.NewEncoder(w)
	return enc.Encode(data)
}

func (JSONCodec) DecodeResponseBody(code uint, r io.ReadCloser, dest interface{}) error {
	defer r.Close()
	if code < 200 {
		return fmt.Errorf("unexpected response status code: %d", code)
	} else if code < 300 {
		dec := json.NewDecoder(r)
		return dec.Decode(dest)
	} else {
		// TODO: handle indexed errors

		var ierr ErrorResponse
		dec := json.NewDecoder(r)
		if err := dec.Decode(&ierr); err != nil {
			return ClientApiErr{fmt.Errorf("failed to decode error response with status code: %d", code)}
		}
		if code < 500 {
			return &InvalidRequest{ierr}
		}
		if code == 503 {
			return &CurrentlySyncing{ierr}
		}
		if code < 600 {
			return &InternalError{ierr}
		}
		return &ierr
	}
}

func (JSONCodec) EncodeRequestBody(w io.Writer, body interface{}) error {
	enc := json.NewEncoder(w) // TODO: different content-types
	return enc.Encode(body)
}

func (JSONCodec) DecodeRequestBody(r io.ReadCloser, dst interface{}) error {
	defer r.Close()
	dec := json.NewDecoder(r)
	return dec.Decode(dst)
}
