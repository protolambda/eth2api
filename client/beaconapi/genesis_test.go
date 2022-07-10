package beaconapi

import (
	"context"
	"testing"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/shared_test"
)

func TestGenesis(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/genesis", "get_genesis",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := Genesis(ctx, cli, new(eth2api.GenesisResponse))
			return err
		})
}
