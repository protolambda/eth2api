package shared_test

import (
	"bytes"
	"context"
	"encoding"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"strings"
	"testing"
)

// All possible API method inputs (params and post-body inputs).
// All of these are optional (not all methods use all inputs).
// A test method can then selectively use them as inputs,
// they still may be nil if they are supposed to be omitted as an optional test-input parameter.
type Input struct {
	ValueStateId      *string  `json:"state_id,omitempty"`
	ValueBlockId      *string  `json:"block_id,omitempty"`
	ValueValidatorId  *string  `json:"validator_id,omitempty"`
	ValueValidatorIds []string `json:"val_ids,omitempty"`

	Slot           *common.Slot              `json:"slot,omitempty"`
	Root           *common.Root              `json:"root,omitempty"`
	Epoch          *common.Epoch             `json:"epoch,omitempty"`
	CommitteeIndex *common.CommitteeIndex    `json:"committee_index,omitempty"`
	ParentRoot     *common.Root              `json:"parent_root,omitempty"`
	Block          *phase0.SignedBeaconBlock `json:"block,omitempty"`
	StatusFilter   []eth2api.ValidatorStatus `json:"validator_statuses,omitempty"`
	// TODO: whole list of possible inputs
}

func (input *Input) StateId() eth2api.StateId {
	if input.ValueStateId == nil {
		return nil
	}
	v, err := eth2api.ParseStateId(*input.ValueStateId)
	if err != nil {
		// not a valid id, but use it anyway to try get the expected error behavior.
		return eth2api.StateIdStrMode(*input.ValueStateId)
	}
	return v
}

func (input *Input) BlockId() eth2api.BlockId {
	if input.ValueBlockId == nil {
		return nil
	}
	v, err := eth2api.ParseBlockId(*input.ValueBlockId)
	if err != nil {
		// not a valid id, but use it anyway to try get the expected error behavior.
		return eth2api.BlockIdStrMode(*input.ValueBlockId)
	}
	return v
}

type mockBadValidatorId string

func (m mockBadValidatorId) ValidatorId() string {
	return string(m)
}

func (input *Input) ValidatorIds() []eth2api.ValidatorId {
	if input.ValueValidatorIds == nil {
		return nil
	}
	ids := input.ValueValidatorIds
	out := make([]eth2api.ValidatorId, len(ids), len(ids))
	var err error
	for i, id := range ids {
		out[i], err = eth2api.ParseValidatorId(id)
		if err != nil {
			// not a valid id, but use it anyway to try get the expected error behavior.
			out[i] = mockBadValidatorId(*input.ValueValidatorId)
		}
	}
	return out
}

func (input *Input) ValidatorId() eth2api.ValidatorId {
	if input.ValueValidatorId == nil {
		return nil
	}
	v, err := eth2api.ParseValidatorId(*input.ValueValidatorId)
	if err != nil {
		// not a valid id, but use it anyway to try get the expected error behavior.
		return mockBadValidatorId(*input.ValueValidatorId)
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
	Description      string       `json:"description"`
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
	return rs.Code, eth2api.JSONCodec{}.DecodeResponseBody(rs.Code,
		ioutil.NopCloser(strings.NewReader(string(rs.Resp))), dest)
}

func (rs requestSpy) Request(_ context.Context, req eth2api.PreparedRequest) eth2api.Response {
	p := "/" + req.Path()
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
					rs.t.Fatalf("failed to encode query key %s: %v", k, err)
				}
				b.Set(k, string(tb))
			} else {
				rs.t.Fatalf("failed to encode query key '%s': unknown type", k)
			}
		}
		p += "?" + b.Encode()
	}
	if rs.ExpectedPath != p {
		rs.t.Fatalf("unexpected request path: '%s', expected: '%s'", p, rs.ExpectedPath)
	}
	method := req.Method()
	if rs.ExpectedPostBody != "" {
		if method != eth2api.POST {
			rs.t.Fatalf("expected POST type, but got enum value: %s", method)
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
				t.Logf("description: %s", c.Description)
				spy := &requestSpy{
					t:            t,
					mockExchange: c,
				}
				// TODO: test cases with timeouts?
				ctx := context.Background()
				err := caseFn(ctx, &c.Input, spy)
				if err != nil {
					if codedErr, ok := err.(eth2api.ApiError); ok {
						if code := codedErr.Code(); code != spy.Code {
							t.Errorf("unexpected code change in bindings: got: %d expected: %d", code, spy.Code)
						}
						// error was expected if code matches.
					} else {
						// e.g. failed to decode response contents
						t.Errorf("unexpected bindings error: %v", err)
					}
				}
			})
		}
	})
}
