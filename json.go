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
	} else if code == 200 {
		dec := json.NewDecoder(r)
		return dec.Decode(dest)
	} else {
		var errMsg ErrorMessage
		dec := json.NewDecoder(r)
		if err := dec.Decode(&errMsg); err != nil {
			return ClientApiErr{fmt.Errorf("failed to decode error response with status code: %d", code)}
		}
		return &errMsg
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
