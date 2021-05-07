package eth2api

import (
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"strconv"
	"strings"
)

type ValidatorId interface {
	ValidatorId() string
}

type ValidatorIdPubkey common.BLSPubkey

func (v ValidatorIdPubkey) ValidatorId() string {
	return common.BLSPubkey(v).String()
}

type ValidatorIdIndex common.ValidatorIndex

func (v ValidatorIdIndex) ValidatorId() string {
	return strconv.FormatUint(uint64(v), 10)
}

func ParseValidatorId(v string) (ValidatorId, error) {
	if strings.HasPrefix(v, "0x") {
		var pub common.BLSPubkey
		if err := pub.UnmarshalText([]byte(v)); err != nil {
			return nil, err
		}
		return ValidatorIdPubkey(pub), nil
	}
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return nil, err
	}
	return ValidatorIdIndex(n), nil
}

type ValidatorIdFilter []ValidatorId

func (vidf ValidatorIdFilter) String() string {
	var out strings.Builder
	for i := range vidf {
		out.WriteString(vidf[i].ValidatorId())
		if i+1 < len(vidf) {
			out.WriteRune(',')
		}
	}
	return out.String()
}

type StateId interface {
	StateId() string
}

type StateIdRoot common.Root

func (r StateIdRoot) StateId() string {
	return common.Root(r).String()
}

type StateIdSlot common.Slot

func (s StateIdSlot) StateId() string {
	return strconv.FormatUint(uint64(s), 10)
}

type StateIdStrMode string

func (m StateIdStrMode) StateId() string {
	return string(m)
}

const (
	StateHead      StateIdStrMode = "head"
	StateGenesis   StateIdStrMode = "genesis"
	StateFinalized StateIdStrMode = "finalized"
	StateJustified StateIdStrMode = "justified"
)

func ParseStateId(v string) (StateId, error) {
	if strings.HasPrefix(v, "0x") {
		var root common.Root
		if err := root.UnmarshalText([]byte(v)); err != nil {
			return nil, err
		}
		return StateIdRoot(root), nil
	}
	asMode := StateIdStrMode(v)
	if asMode == StateHead || asMode == StateGenesis || asMode == StateFinalized || asMode == StateJustified {
		return asMode, nil
	}
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return nil, err
	}
	return StateIdSlot(n), nil
}

// Block identifier
// Can have different shapes. E.g. BlockHead, BlockIdSlot(123), BlockIdRoot(Root{0x...})
type BlockId interface {
	BlockId() string
}

type BlockIdRoot common.Root

func (r BlockIdRoot) BlockId() string {
	return common.Root(r).String()
}

type BlockIdSlot common.Slot

func (s BlockIdSlot) BlockId() string {
	return strconv.FormatUint(uint64(s), 10)
}

type BlockIdStrMode string

func (m BlockIdStrMode) BlockId() string {
	return string(m)
}

const (
	BlockHead      BlockIdStrMode = "head"
	BlockGenesis   BlockIdStrMode = "genesis"
	BlockFinalized BlockIdStrMode = "finalized"
)

func ParseBlockId(v string) (BlockId, error) {
	if strings.HasPrefix(v, "0x") {
		var root common.Root
		if err := root.UnmarshalText([]byte(v)); err != nil {
			return nil, err
		}
		return BlockIdRoot(root), nil
	}
	asMode := BlockIdStrMode(v)
	if asMode == BlockHead || asMode == BlockGenesis || asMode == BlockFinalized {
		return asMode, nil
	}
	n, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return nil, err
	}
	return BlockIdSlot(n), nil
}

type StatusFilter []ValidatorStatus

func (sf StatusFilter) String() string {
	var out strings.Builder
	for i := range sf {
		out.WriteString(string(sf[i]))
		if i+1 < len(sf) {
			out.WriteRune(',')
		}
	}
	return out.String()
}
