package shared_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/protolambda/eth2api"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
)

type Input struct {
	ValueStateId string `json:"state_id,omitempty"`
	ValueBlockId string `json:"block_id,omitempty"`
	// TODO: whole list of possible inputs
}

func (input *Input) StateId() eth2api.StateId {
	v, err := eth2api.ParseStateId(input.ValueStateId)
	if err != nil {
		panic(err) // invalid test resource
	}
	return v
}

func (input *Input) BlockId() eth2api.BlockId {
	v, err := eth2api.ParseBlockId(input.ValueBlockId)
	if err != nil {
		panic(err) // invalid test resource
	}
	return v
}

func MustExist(exists bool, err error) error {
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("expected it to exist")
	}
	return nil
}

type mockResponse string

func (mr *mockResponse) UnmarshalJSON(v []byte) error {
	// lazy json decoding
	*mr = mockResponse(v)
	return nil
}

type mockExchange struct {
	Input            Input        `json:"input"`
	ExpectedPath     string       `json:"path"`
	ExpectedPostBody string       `json:"post_body,omitempty"`
	Code             uint         `json:"code"`
	Resp             mockResponse `json:"response"`
}

type requestSpy struct {
	t *testing.T
	*mockExchange
}

func (rs requestSpy) Decode(dest interface{}) (code uint, err error) {
	return eth2api.DecodeBody(rs.Code, ioutil.NopCloser(strings.NewReader(string(rs.Resp))), dest)
}

func (rs requestSpy) Request(ctx context.Context, req eth2api.Request) eth2api.Response {
	p := "/" + req.Path()
	if q := req.Query(); q != nil {
		var b url.Values
		for k, v := range req.Query() {
			b.Set(k, fmt.Sprintf("%s", v))
		}
		p += "?" + b.Encode()
	}
	if rs.ExpectedPath != p {
		rs.t.Fatalf("unexpected request path: '%s', expected: '%s'", p, rs.ExpectedPath)
	}
	method := req.Method()
	if rs.ExpectedPostBody != "" {
		if method != eth2api.POST {
			rs.t.Fatalf("expected POST type (enum %d), but got enum value: %d", eth2api.POST, method)
		}
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf) // TODO: different content-types
		if err := enc.Encode(req.Body()); err != nil {
			rs.t.Fatalf("failed to encode POST body: %v", err)
		}
		if got := buf.String(); got != rs.ExpectedPostBody {
			rs.t.Fatalf("unexpected POST body contents:\ngot:\n---\n%s\n---\nexpected:\n---\n%s\n---\n", got, rs.ExpectedPostBody)
		}
	}
	return rs
}

// Loads one or more test cases.
// Each test case can be executed as an artificial API client which checks what it is being called with,
// and then outputs the test mock response.
func loadTests(t *testing.T, sourcePath string) []*mockExchange {
	f, err := os.Open(sourcePath)
	if err != nil {
		t.Fatalf("failed to open test source: %v", err)
		return nil
	}
	defer f.Close()
	var mocks []*mockExchange
	if err := json.NewDecoder(f).Decode(&mocks); err != nil {
		t.Fatalf("failed to load test source: %v", err)
	}
	return mocks
}

func RunAll(t *testing.T, testsDir string, name string, caseFn func(ctx context.Context, input *Input, cli eth2api.Client) error) {
	t.Run(name, func(t *testing.T) {
		cases := loadTests(t, path.Join(testsDir, name+".json"))
		for i, c := range cases {
			t.Run(fmt.Sprintf("case_%d", i), func(t *testing.T) {
				spy := &requestSpy{
					t:            t,
					mockExchange: c,
				}
				// TODO: test cases with timeouts?
				ctx := context.Background()
				err := caseFn(ctx, &c.Input, spy)
				if err != nil {
					t.Error(err)
				}
			})
		}
	})
}
