package beaconapi

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/merge"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"github.com/protolambda/zrnt/eth2/beacon/sharding"
)

// Serves attestations included in requested block.
func BlockAttestations(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET,
		"/eth/v1/beacon/blocks/:blockId/attestations",
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
			block, ok := blockEnvelop.SignedBlock.(*phase0.SignedBeaconBlock)
			if !ok {
				return eth2api.RespondInternalError(fmt.Errorf(
					"only supporting attestation fetching of phase0 blocks, but got %T", blockEnvelop.SignedBlock))
			}
			return eth2api.RespondOK(eth2api.Wrap(block.Message.Body.Attestations))
		})
}

// Serves phase0 block details for given block id.
func Block(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "/eth/v1/beacon/blocks/:blockId",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			blockId, err := eth2api.ParseBlockId(req.Param("blockId"))
			if err != nil {
				return eth2api.RespondBadInput(err)
			}
			entry, ok := backend.BlockLookup(blockId)
			if !ok {
				return eth2api.RespondNotFound("Block not found")
			}
			blockEnv, err := backend.BlockDB.Get(ctx, entry.BlockRoot())
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block: %v", err))
			}
			if blockEnv == nil {
				return eth2api.RespondNotFound("Block not found")
			}
			return eth2api.RespondOK(eth2api.Wrap(blockEnv.SignedBlock))
		})
}

// Serves versioned block details for given block id.
func Blockv2(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "/eth/v2/beacon/blocks/:blockId",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			blockId, err := eth2api.ParseBlockId(req.Param("blockId"))
			if err != nil {
				return eth2api.RespondBadInput(err)
			}
			entry, ok := backend.BlockLookup(blockId)
			if !ok {
				return eth2api.RespondNotFound("Block not found")
			}
			blockEnv, err := backend.BlockDB.Get(ctx, entry.BlockRoot())
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block: %v", err))
			}
			if blockEnv == nil {
				return eth2api.RespondNotFound("Block not found")
			}
			var version string
			switch blockEnv.SignedBlock.(type) {
			case *phase0.SignedBeaconBlock:
				version = "phase0"
			case *altair.SignedBeaconBlock:
				version = "altair"
			case *merge.SignedBeaconBlock:
				version = "merge"
			case *sharding.SignedBeaconBlock:
				version = "sharding"
			default:
				return eth2api.RespondInternalError(fmt.Errorf("unknown block type %T", blockEnv.SignedBlock))
			}
			return eth2api.RespondOK(&eth2api.VersionedBeaconBlock{Version: version, Data: blockEnv.SignedBlock})
		})
}

type slotDecodeLookahead struct {
	Message struct {
		Slot common.Slot `json:"slot"`
	} `json:"message"`
}

type slotHack struct {
	backend *BeaconBackend
	dest    *common.BeaconBlockEnvelope
}

func (h *slotHack) UnmarshalJSON(b []byte) error {
	var slotData slotDecodeLookahead
	if err := json.Unmarshal(b, slotData); err != nil {
		return err
	}

	forkDigest := common.ComputeForkDigest(
		h.backend.ForkDecoder.Spec.ForkVersion(slotData.Message.Slot),
		h.backend.Chain.Genesis().ValidatorsRoot)

	dest, err := h.backend.ForkDecoder.AllocBlock(forkDigest)
	if err != nil {
		return fmt.Errorf("unrecognized fork: %v", err)
	}

	if err := json.Unmarshal(b, dest); err != nil {
		return err
	}

	h.dest = dest.Envelope(h.backend.Spec, forkDigest)
	return nil
}

// Instructs the beacon node to broadcast a newly signed beacon block to the beacon network,
// to be included in the beacon chain. The beacon node is not required to validate the signed `BeaconBlock`,
// and a successful response (20X, i.e. no error returned) only indicates that the broadcast has been successful.
// The beacon node is expected to integrate the new block into its state, and therefore validate the block internally,
// however blocks which fail the validation are still broadcast but a different status code is returned (202).
func PublishBlock(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.POST, "/eth/v1/beacon/blocks",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var block slotHack
			if err := req.DecodeBody(&block); err != nil {
				return eth2api.RespondBadInput(err)
			}
			blockEnvelop := block.dest
			syncing, err := backend.Publisher.PublishBlock(ctx, blockEnvelop)

			// handle even if we cannot publish it, to keep liveness in case sync is bad.
			if _, err2 := backend.BlockDB.Store(ctx, blockEnvelop); err2 != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to store block: %v", err2))
			}
			if err2 := backend.Chain.AddBlock(ctx, blockEnvelop); err2 != nil {
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
	return eth2api.MakeRoute(eth2api.GET, "/eth/v1/beacon/blocks/:blockId",
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
