package beaconapi

import (
	"context"
	"fmt"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"strconv"
)

// Serve block header for given block id.
func BlockHeader(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/headers/:blockId",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			blockId, err := eth2api.ParseBlockId(req.Param("blockId"))
			if err != nil {
				return eth2api.RespondBadInput(err)
			}
			entry, ok := backend.BlockLookup(blockId)
			if !ok {
				return eth2api.RespondNotFound("Block not found")
			}
			blockEnvelop, err := backend.BlockDB.Get(ctx, entry.BlockRoot())
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block: %v", err))
			}
			if blockEnvelop == nil {
				return eth2api.RespondNotFound("Block not found")
			}
			canon, ok := backend.Chain.ByCanonStep(entry.Step())
			if !ok {
				return eth2api.RespondInternalError(fmt.Errorf("failed to determine if entry is canonical"))
			}
			// TODO: support alternative forks
			block, ok := blockEnvelop.SignedBlock.(*phase0.SignedBeaconBlock)
			if !ok {
				return eth2api.RespondInternalError(fmt.Errorf("only supporting phase0 block headers currently"))
			}
			header := block.SignedHeader(backend.Spec)
			out := eth2api.BeaconBlockHeaderAndInfo{
				Root:      entry.BlockRoot(),
				Canonical: canon.BlockRoot() == entry.BlockRoot(),
				Header:    *header,
			}
			return eth2api.RespondOK(eth2api.Wrap(&out))
		})
}

// Serves block headers matching given query. By default it will serve current head slot blocks.
func BlockHeaders(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/headers",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var parentRootFilter *common.Root
			parentRootVals, ok := req.Query("parent_root")
			if ok && len(parentRootVals) > 0 {
				var root common.Root
				if err := root.UnmarshalText([]byte(parentRootVals[0])); err != nil {
					return eth2api.RespondBadInput(fmt.Errorf("bad parent_root query param: %v", err))
				}
				parentRootFilter = &root
			}
			var slotFilter *common.Slot
			slotVals, ok := req.Query("slot")
			if ok && len(slotVals) > 0 {
				v, err := strconv.ParseUint(slotVals[0], 0, 64)
				if err != nil {
					return eth2api.RespondBadInput(fmt.Errorf("bad slot query param: %v", err))
				}
				slot := common.Slot(v)
				slotFilter = &slot
			}
			results, err := backend.Chain.Search(parentRootFilter, slotFilter)
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to search for headers: %v", err))
			}
			data := make([]eth2api.BeaconBlockHeaderAndInfo, 0, len(results))
			for _, res := range results {
				blockEnvelop, err := backend.BlockDB.Get(ctx, res.BlockRoot())
				if err != nil {
					return eth2api.RespondInternalError(fmt.Errorf("failed to load block: %v", err))
				}
				if blockEnvelop == nil {
					return eth2api.RespondNotFound("Block not found")
				}
				block, ok := blockEnvelop.SignedBlock.(*phase0.SignedBeaconBlock)
				if !ok {
					return eth2api.RespondInternalError(fmt.Errorf("only supporting phase0 block headers currently"))
				}
				data = append(data, eth2api.BeaconBlockHeaderAndInfo{
					Root:      res.BlockRoot(),
					Canonical: res.Canonical,
					Header:    *block.SignedHeader(backend.Spec),
				})
			}
			return eth2api.RespondOK(eth2api.Wrap(data))
		})
}
