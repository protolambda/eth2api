package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	chain2 "github.com/protolambda/zrnt/eth2/chain"
	"github.com/protolambda/zrnt/eth2/db/blocks"
	"github.com/protolambda/zrnt/eth2/pool"
)

type Publisher interface {
	PublishBlock(ctx context.Context, block *phase0.SignedBeaconBlock) (syncing bool, err error)
	PublishAttestation(ctx context.Context, att *phase0.Attestation) (err error)
	PublishAttesterSlashing(ctx context.Context, att *phase0.AttesterSlashing) (err error)
	PublishProposerSlashing(ctx context.Context, att *phase0.ProposerSlashing) (err error)
	PublishVoluntaryExits(ctx context.Context, att *phase0.SignedVoluntaryExit) (err error)
}

type BeaconBackend struct {
	Spec      *common.Spec
	Chain     chain2.FullChain
	BlockDB   blocks.DB
	Publisher Publisher

	// TODO move pools to interface
	AttestationPool      *pool.AttestationPool
	AttesterSlashingPool *pool.AttesterSlashingPool
	ProposerSlashingPool *pool.ProposerSlashingPool
	VoluntaryExitPool    *pool.VoluntaryExitPool
}

func (backend *BeaconBackend) BlockLookup(blockId eth2api.BlockId) (entry chain2.ChainEntry, ok bool) {
	switch id := blockId.(type) {
	case eth2api.BlockIdRoot:
		return backend.Chain.ByBlock(common.Root(id))
	case eth2api.BlockIdSlot:
		// prefer a slot entry that includes the block.
		entry, ok = backend.Chain.ByCanonStep(chain2.AsStep(common.Slot(id), true))
		if !ok {
			entry, ok = backend.Chain.ByCanonStep(chain2.AsStep(common.Slot(id), false))
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
			return backend.Chain.ByCanonStep(chain2.AsStep(common.Slot(0), true))
		default:
			return nil, false
		}
	default:
		return nil, false
	}
}
