package beaconapi

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/bellatrix"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
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
			blockRoot, err := entry.BlockRoot()
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block root: %v", err))
			}
			blockEnvelop, err := backend.BlockDB.Get(entry.Step().Slot(), blockRoot)
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block: %v", err))
			}
			if blockEnvelop == nil {
				return eth2api.RespondNotFound("Block not found")
			}
			var atts phase0.Attestations
			switch x := blockEnvelop.Body.(type) {
			case *phase0.BeaconBlockBody:
				atts = x.Attestations
			case *altair.BeaconBlockBody:
				atts = x.Attestations
			default:
				return eth2api.RespondInternalError(fmt.Errorf("unrecongized beacon block body type: %T", x))
			}
			return eth2api.RespondOK(eth2api.Wrap(atts))
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
			blockRoot, err := entry.BlockRoot()
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block root: %v", err))
			}
			blockEnv, err := backend.BlockDB.Get(entry.Step().Slot(), blockRoot)
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block: %v", err))
			}
			if blockEnv == nil {
				return eth2api.RespondNotFound("Block not found")
			}
			data, err := beacon.EnvelopeToSignedBeaconBlock(blockEnv)
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to construct typed beacon block: %v", err))
			}
			return eth2api.RespondOK(eth2api.Wrap(data))
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
			blockRoot, err := entry.BlockRoot()
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block root: %v", err))
			}
			blockEnv, err := backend.BlockDB.Get(entry.Step().Slot(), blockRoot)
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block: %v", err))
			}
			if blockEnv == nil {
				return eth2api.RespondNotFound("Block not found")
			}
			data, err := beacon.EnvelopeToSignedBeaconBlock(blockEnv)
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to construct typed beacon block: %v", err))
			}
			var version string
			switch data.(type) {
			case *phase0.SignedBeaconBlock:
				version = "phase0"
			case *altair.SignedBeaconBlock:
				version = "altair"
			case *bellatrix.SignedBeaconBlock:
				version = "bellatrix"
			default:
				return eth2api.RespondInternalError(fmt.Errorf("unknown block type %T", data))
			}
			return eth2api.RespondOK(&eth2api.VersionedBeaconBlock{Version: version, Data: data})
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
	if err := json.Unmarshal(b, &slotData); err != nil {
		return err
	}

	forkDigest := common.ComputeForkDigest(
		h.backend.ForkDecoder.Spec.ForkVersion(slotData.Message.Slot),
		h.backend.Chain.Genesis().ValidatorsRoot)

	alloc, err := h.backend.ForkDecoder.BlockAllocator(forkDigest)
	if err != nil {
		return fmt.Errorf("unrecognized fork: %v", err)
	}
	dest := alloc()

	if err := json.Unmarshal(b, dest); err != nil {
		return err
	}

	h.dest = dest.Envelope(h.backend.Spec, forkDigest)
	return nil
}

// Instructs the beacon node to broadcast a newly signed beacon block to the beacon network, to be included in the beacon chain.
// The beacon node is not required to validate the signed BeaconBlock, and a successful response (20X) only indicates that the broadcast has been successful.
// The beacon node is expected to integrate the new block into its state, and therefore validate the block internally,
// however blocks which fail the validation are still broadcast but a different status code is returned (202)
func PublishBlock(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.POST, "/eth/v1/beacon/blocks",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			var block slotHack
			if err := req.DecodeBody(&block); err != nil {
				return eth2api.RespondBadInput(err)
			}
			blockEnvelop := block.dest
			syncing, err := backend.Publisher.PublishBlock(ctx, blockEnvelop)
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to publish block: %v", err))
			} else if syncing {
				return eth2api.RespondSyncing("beacon is syncing, added it to the chain, but cannot publish block")
			}

			if err := backend.ProcessBlock(ctx, blockEnvelop); err != nil {
				return eth2api.RespondAccepted(fmt.Errorf("published block, but failed to process locally: %v", err))
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
			blockRoot, err := entry.BlockRoot()
			if err != nil {
				return eth2api.RespondInternalError(fmt.Errorf("failed to load block root: %v", err))
			}
			out := eth2api.RootResponse{Root: blockRoot}
			return eth2api.RespondOK(eth2api.Wrap(&out))
		})
}
