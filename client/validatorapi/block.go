package validatorapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// Requests a beacon node to produce a valid block, which can then be signed by a validator.
//
// Err will be non-nil when syncing.
func ProduceBlock(ctx context.Context, cli eth2api.Client,
	slot common.Slot, randaoReveal common.BLSSignature, graffiti *common.Root, dest *phase0.BeaconBlock) (syncing bool, err error) {
	q := eth2api.Query{
		"randao_reveal": randaoReveal,
	}
	if graffiti != nil {
		q["graffiti"] = graffiti
	}
	req := eth2api.FmtQueryGET(q, "eth/v1/validator/blocks/%d", slot)
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(eth2api.Wrap(dest))
	syncing = code == 503
	return
}
