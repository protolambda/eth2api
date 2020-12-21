package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/shared_test"
	"testing"
)

func TestBlockHeader(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/headers", "get_header",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := BlockHeader(ctx, cli, input.BlockId(), new(eth2api.BeaconBlockHeaderAndInfo))
			return err
		})
}

func TestBlockHeaders(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/headers", "get_headers",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := BlockHeaders(ctx, cli, input.Slot, input.ParentRoot, new([]eth2api.BeaconBlockHeaderAndInfo))
			return err
		})
}
