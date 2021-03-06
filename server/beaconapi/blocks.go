package beaconapi

import (
	"context"
	"fmt"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/db/blocks"
)

// Serves attestations included in requested block.
func BlockAttestations(backend *BeaconBackend) eth2api.Route {
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

// Serves block details for given block id.
func Block(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/blocks/:blockId",
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
			return eth2api.RespondOK(eth2api.Wrap(&block))
		})
}

// Instructs the beacon node to broadcast a newly signed beacon block to the beacon network,
// to be included in the beacon chain. The beacon node is not required to validate the signed `BeaconBlock`,
// and a successful response (20X, i.e. no error returned) only indicates that the broadcast has been successful.
// The beacon node is expected to integrate the new block into its state, and therefore validate the block internally,
// however blocks which fail the validation are still broadcast but a different status code is returned (202).
func PublishBlock(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.POST, "eth/v1/beacon/blocks",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var block beacon.SignedBeaconBlock
			if err := req.DecodeBody(&block); err != nil {
				return eth2api.RespondBadInput(err)
			}

			syncing, err := backend.Publisher.PublishBlock(ctx, &block)

			// handle even if we cannot publish it, to keep liveness in case sync is bad.
			if _, err2 := backend.BlockDB.Store(ctx, blocks.WithRoot(backend.Spec, &block)); err2 != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to store block: %v", err2))
			}
			if err2 := backend.Chain.AddBlock(ctx, &block); err2 != nil {
				return eth2api.RespondBadInput(fmt.Errorf("failed to process block: %v", err2))
			}

			if err != nil {
				return eth2api.RespondInternalError(err)
			} else if syncing {
				return eth2api.RespondSyncing("beacon is syncing, added it to the chain, but cannot publish block")
			}

			return eth2api.RespondOKMsg("processed and published block")
		})
}

// Serves hashTreeRoot of BeaconBlock/BeaconBlockHeader.
func BlockRoot(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/blocks/:blockId",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			blockId, err := eth2api.ParseBlockId(req.Param("blockId"))
			if err != nil {
				return eth2api.RespondBadInput(err)
			}
			entry, ok := backend.BlockLookup(blockId)
			if !ok {
				return eth2api.RespondNotFound("Block not found")
			}
			out := eth2api.RootResponse{Root: entry.BlockRoot()}
			return eth2api.RespondOK(eth2api.Wrap(&out))
		})
}
