package configapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
)

// Retrieve specification configuration used on this node.
func Spec(ctx context.Context, cli eth2api.Client, dest *beacon.Phase0Config) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("eth/v1/config/spec"), dest)
}
