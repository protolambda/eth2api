package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Retrieve details of the chain's genesis which can be used to identify chain.
func Genesis(ctx context.Context, cli eth2api.Client, dest *eth2api.GenesisResponse) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.PlainGET("eth/v1/beacon/genesis"), eth2api.Wrap(dest))
}
