package nodeapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Requests the beacon node to describe if it's currently syncing or not, and if it is, what block it is up to.
func SyncingStatus(ctx context.Context, cli eth2api.Client, dest *eth2api.SyncingStatus) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("eth/v1/node/syncing"), dest)
}
