package debugapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
)

// Returns full BeaconState object for given stateId.
func State(ctx context.Context, cli eth2api.Client, stateId eth2api.StateId, dest *beacon.BeaconState) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/debug/state/%s", stateId.StateId()), eth2api.Wrap(dest))
}
