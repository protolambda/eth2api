package configapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
)

// Retrieve all scheduled upcoming forks this node is aware of.
func ForkSchedule(ctx context.Context, cli eth2api.Client, dest *[]beacon.Fork) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("eth/v1/config/fork_schedule"), dest)
}
