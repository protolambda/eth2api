package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
)

// Retrieves attestations included in requested block.
func BlockAttestations(ctx context.Context, cli eth2api.Client, blockId eth2api.BlockId, dest *[]beacon.Attestation) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/blocks/%s/attestations", blockId.BlockId()), dest)
}

// Retrieves block details for given block id.
func Block(ctx context.Context, cli eth2api.Client, blockId eth2api.BlockId, dest *beacon.SignedBeaconBlock) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/blocks/%s", blockId.BlockId()), dest)
}

// Instructs the beacon node to broadcast a newly signed beacon block to the beacon network,
// to be included in the beacon chain. The beacon node is not required to validate the signed `BeaconBlock`,
// and a successful response (20X, i.e. no error returned) only indicates that the broadcast has been successful.
// The beacon node is expected to integrate the new block into its state, and therefore validate the block internally,
// however blocks which fail the validation are still broadcast but a different status code is returned
// (202, `valid` will be false)
func PublishBlock(ctx context.Context, cli eth2api.Client, block *beacon.SignedBeaconBlock) (valid bool, err error) {
	req := eth2api.BodyPOST("eth/v1/beacon/blocks", block)
	resp := cli.Request(ctx, req)
	if err := resp.Err(); err != nil {
		if err.Code() == 202 {
			return false, nil
		}
		return false, err
	}
	valid = true
	return
}

// Retrieves block header for given block id.
func BlockHeader(ctx context.Context, cli eth2api.Client, blockId eth2api.BlockId, dest *beacon.SignedBeaconBlockHeader) (exists bool, err error) {
	if blockId == nil {
		return false, eth2api.MissingRequiredParamErr
	}
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/blocks/%s", blockId.BlockId()), dest)
}

// Retrieves block headers matching given query. By default it will fetch current head slot blocks.
func BlockHeaders(ctx context.Context, cli eth2api.Client, slot *beacon.Slot, parentRoot *beacon.Root, dest *[]beacon.SignedBeaconBlockHeader) (exists bool, err error) {
	var q eth2api.Query
	if slot != nil {
		if parentRoot != nil {
			q = eth2api.Query{"slot": *slot, "parent_root": *parentRoot}
		} else {
			q = eth2api.Query{"slot": *slot}
		}
	} else if parentRoot != nil {
		q = eth2api.Query{"parent_root": *parentRoot}
	}
	return eth2api.SimpleRequest(ctx, cli, eth2api.QueryGET(q, "eth/v1/beacon/blocks/header"), dest)
}

// Retrieves hashTreeRoot of BeaconBlock/BeaconBlockHeader.
func BlockRoot(ctx context.Context, cli eth2api.Client, blockId eth2api.BlockId, dest *beacon.Root) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/blocks/%s/root", blockId.BlockId()), dest)
}
