package eth2api

import (
	"github.com/protolambda/zrnt/eth2/beacon"
	"strconv"
)

type StateId interface {
	StateId() string
}

type StateIdRoot beacon.Root

func (r StateIdRoot) StateId() string {
	return beacon.Root(r).String()
}

type StateIdSlot beacon.Slot

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

// Block identifier
// Can have different shapes. E.g. BlockHead, BlockIdSlot(123), BlockIdRoot(Root{0x...})
type BlockId interface {
	BlockId() string
}

type BlockIdRoot beacon.Root

func (r BlockIdRoot) BlockId() string {
	return beacon.Root(r).String()
}

type BlockIdSlot beacon.Slot

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
