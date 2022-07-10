package nodeapi

import (
	"context"

	"github.com/protolambda/eth2api"
)

// Checks node health. Healthy = no error. May be syncing, capable of serving incomplete data.
//
// Err will be non-nil when syncing.
func Health(ctx context.Context, cli eth2api.Client) (syncing bool, err error) {
	resp := cli.Request(ctx, eth2api.PlainGET("/eth/v1/node/health"))
	var code uint
	code, err = resp.Decode(nil)
	syncing = code == 206
	return
}
