package debugapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// Retrieves phase0 BeaconState object for given stateId.
func BeaconState(ctx context.Context, cli eth2api.Client, stateId eth2api.StateId, dest *phase0.BeaconState) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("/eth/v1/debug/beacon/states/%s", stateId.StateId()), eth2api.Wrap(dest))
}

// Retrieves versioned BeaconState object for given stateId.
func BeaconStateV2(ctx context.Context, cli eth2api.Client, stateId eth2api.StateId, dest *eth2api.VersionedBeaconState) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("/eth/v2/debug/beacon/states/%s", stateId.StateId()), dest)
}
