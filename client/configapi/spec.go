package configapi

import (
	"context"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

// Retrieve specification configuration used on this node.
func Spec(ctx context.Context, cli eth2api.Client, dest *common.Spec) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("/eth/v1/config/spec"), eth2api.Wrap(dest))
}
