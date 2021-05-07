package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// Retrieves attestations included in requested block.
func BlockAttestations(ctx context.Context, cli eth2api.Client, blockId eth2api.BlockId, dest *[]phase0.Attestation) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/blocks/%s/attestations", blockId.BlockId()), eth2api.Wrap(dest))
}

// Retrieves block details for given block id.
func Block(ctx context.Context, cli eth2api.Client, blockId eth2api.BlockId, dest *phase0.SignedBeaconBlock) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/blocks/%s", blockId.BlockId()), eth2api.Wrap(dest))
}

// Instructs the beacon node to broadcast a newly signed beacon block to the beacon network,
// to be included in the beacon chain. The beacon node is not required to validate the signed `BeaconBlock`,
// and a successful response (20X, i.e. no error returned) only indicates that the broadcast has been successful.
// The beacon node is expected to integrate the new block into its state, and therefore validate the block internally,
// however blocks which fail the validation are still broadcast but a different status code is returned
// (202, `valid` will be false)
func PublishBlock(ctx context.Context, cli eth2api.Client, block *phase0.SignedBeaconBlock) (valid bool, err error) {
	req := eth2api.BodyPOST("eth/v1/beacon/blocks", block)
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(nil)
	valid = code != 202
	return
}

// Retrieves hashTreeRoot of BeaconBlock/BeaconBlockHeader.
func BlockRoot(ctx context.Context, cli eth2api.Client, blockId eth2api.BlockId) (root common.Root, exists bool, err error) {
	var dest eth2api.RootResponse
	exists, err = eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/blocks/%s/root", blockId.BlockId()), eth2api.Wrap(&dest))
	root = dest.Root
	return
}
