package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/pool"
)

type Publisher interface {
	PublishBlock(ctx context.Context, block *common.BeaconBlockEnvelope) (syncing bool, err error)
	PublishAttestation(ctx context.Context, att *phase0.Attestation) (err error)
	PublishAttesterSlashing(ctx context.Context, sl *phase0.AttesterSlashing) (err error)
	PublishProposerSlashing(ctx context.Context, sl *phase0.ProposerSlashing) (err error)
	PublishVoluntaryExit(ctx context.Context, exit *phase0.SignedVoluntaryExit) (err error)
	PublishSyncCommitteeMessage(ctx context.Context, msg *altair.SyncCommitteeMessage) error
}

type BlockReader interface {
	Get(slot common.Slot, root common.Root) (*common.BeaconBlockEnvelope, error)
}

type AttestationPool interface {
	Search(opts ...pool.AttSearchOption) (out []*phase0.Attestation)
	AddAttestation(ctx context.Context, att *phase0.Attestation, committee common.CommitteeIndices) error
}

type AttesterSlashingPool interface {
	All() []*phase0.AttesterSlashing
	AddAttesterSlashing(ctx context.Context, sl *phase0.AttesterSlashing) error
}

type ProposerSlashingPool interface {
	All() []*phase0.ProposerSlashing
	AddProposerSlashing(ctx context.Context, sl *phase0.ProposerSlashing) error
}

type VoluntaryExitPool interface {
	All() []*phase0.SignedVoluntaryExit
	AddVoluntaryExit(ctx context.Context, exit *phase0.SignedVoluntaryExit) error
}

type SyncCommitteePool interface {
	AddSyncCommitteeMessage(ctx context.Context, msg *altair.SyncCommitteeMessage) error
}

type BeaconBackend struct {
	Spec      *common.Spec
	Chain     beacon.Chain
	BlockDB   BlockReader
	Publisher Publisher

	ProcessBlock func(ctx context.Context, block *common.BeaconBlockEnvelope) error

	ForkDecoder *beacon.ForkDecoder

	AttestationPool      AttestationPool
	AttesterSlashingPool AttesterSlashingPool
	ProposerSlashingPool ProposerSlashingPool
	VoluntaryExitPool    VoluntaryExitPool
	SyncCommitteePool    SyncCommitteePool
}

func (backend *BeaconBackend) BlockLookup(blockId eth2api.BlockId) (entry beacon.ChainEntry, ok bool) {
	switch id := blockId.(type) {
	case eth2api.BlockIdRoot:
		return backend.Chain.ByBlock(common.Root(id))
	case eth2api.BlockIdSlot:
		// prefer a slot entry that includes the block.
		entry, ok = backend.Chain.ByCanonStep(common.AsStep(common.Slot(id), true))
		if !ok {
			entry, ok = backend.Chain.ByCanonStep(common.AsStep(common.Slot(id), false))
		}
		return
	case eth2api.BlockIdStrMode:
		switch id {
		case eth2api.BlockHead:
			entry, err := backend.Chain.Head()
			return entry, err == nil
		case eth2api.BlockFinalized:
			entry, err := backend.Chain.Finalized()
			return entry, err == nil
		case eth2api.BlockGenesis:
			return backend.Chain.ByCanonStep(common.AsStep(common.Slot(0), true))
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}
