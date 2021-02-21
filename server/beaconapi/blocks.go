package beaconapi

import (
	"context"
	"fmt"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/chain"
	"github.com/protolambda/zrnt/eth2/db/blocks"
)

type BeaconBackend struct {
	Chain   chain.FullChain
	BlockDB blocks.DB
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

func BlockAttestations(backend BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET,
		"eth/v1/beacon/blocks/:blockId/attestations",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			blockId, err := eth2api.ParseBlockId(req.Param("blockId"))
			if err != nil {
				return eth2api.RespondBadInput(err)
			}
			entry, ok := backend.BlockLookup(blockId)
			if !ok {
				return eth2api.RespondNotFound("Block not found")
			}
			var block beacon.SignedBeaconBlock
			exists, err := backend.BlockDB.Get(ctx, entry.BlockRoot(), &block)
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block: %v", err))
			}
			if !exists {
				return eth2api.RespondNotFound("Block not found")
			}
			return eth2api.RespondOK(eth2api.Wrap(block.Message.Body.Attestations))
		})
}
