package debugapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Retrieves all possible chain heads (leaves of fork choice tree).
func BeaconChainHeads(ctx context.Context, cli eth2api.Client, dest *[]eth2api.ChainHead) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("/eth/v1/debug/beacon/heads"), eth2api.Wrap(dest))
}
