package nodeapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Retrieves data about the node's network presence.
func NetworkIdentity(ctx context.Context, cli eth2api.Client, dest *eth2api.NetworkIdentity) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("/eth/v1/node/network_identity"), eth2api.Wrap(dest))
}
