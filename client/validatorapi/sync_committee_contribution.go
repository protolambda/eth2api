package validatorapi

import (
	"context"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/view"
)

// Requests that the beacon node produces a sync committee contribution.
//
// Err will be non-nil when syncing.
func ProduceSyncCommitteeContribution(ctx context.Context, cli eth2api.Client, slot common.Slot, subCommitteeIndex uint64,
	beaconBlockRoot common.Root, dest *altair.SyncCommitteeContribution) (syncing bool, err error) {
	q := eth2api.Query{
		"slot":               slot,
		"subcommittee_index": view.Uint64View(subCommitteeIndex),
		"beacon_block_root":  beaconBlockRoot,
	}
	req := eth2api.QueryGET(q, "/eth/v1/validator/sync_committee_contribution")
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(eth2api.Wrap(dest))
	syncing = code == 503
	return
}
