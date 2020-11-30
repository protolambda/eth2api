package nodeapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Checks node health. Healthy = no error. May be syncing, capable of serving incomplete data.
func ChainHeads(ctx context.Context, cli eth2api.Client) (syncing bool, err error) {
	resp := cli.Request(ctx, eth2api.PlainGET("eth/v1/node/health"))
	if err := resp.Err(); err != nil {
		if err.Code() == 206 {
			return true, nil
		}
		return false, err
	}
	return false, nil
}
