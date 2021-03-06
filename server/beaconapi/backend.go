package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/chain"
	"github.com/protolambda/zrnt/eth2/db/blocks"
	"github.com/protolambda/zrnt/eth2/pool"
)

type Publisher interface {
	PublishBlock(ctx context.Context, block *beacon.SignedBeaconBlock) (syncing bool, err error)
	PublishAttestation(ctx context.Context, att *beacon.Attestation) (err error)
	PublishAttesterSlashing(ctx context.Context, att *beacon.AttesterSlashing) (err error)
	PublishProposerSlashing(ctx context.Context, att *beacon.ProposerSlashing) (err error)
	PublishVoluntaryExits(ctx context.Context, att *beacon.SignedVoluntaryExit) (err error)
}

type BeaconBackend struct {
	Spec            *beacon.Spec
	Chain           chain.FullChain
	BlockDB         blocks.DB
	Publisher       Publisher

	// TODO move pools to interface
	AttestationPool *pool.AttestationPool
	AttesterSlashingPool *pool.AttesterSlashingPool
	ProposerSlashingPool *pool.ProposerSlashingPool
	VoluntaryExitPool *pool.VoluntaryExitPool
}

func (backend *BeaconBackend) BlockLookup(blockId eth2api.BlockId) (entry chain.ChainEntry, ok bool) {
	switch id := blockId.(type) {
	case eth2api.BlockIdRoot:
		return backend.Chain.ByBlock(beacon.Root(id))
	case eth2api.BlockIdSlot:
		// prefer a slot entry that includes the block.
		entry, ok = backend.Chain.ByCanonStep(chain.AsStep(beacon.Slot(id), true))
		if !ok {
			entry, ok = backend.Chain.ByCanonStep(chain.AsStep(beacon.Slot(id), false))
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
			return backend.Chain.ByCanonStep(chain.AsStep(beacon.Slot(0), true))
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}
