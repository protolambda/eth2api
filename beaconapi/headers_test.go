package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/shared_test"
	"github.com/protolambda/zrnt/eth2/beacon"
	"testing"
)

func TestBlockHeader(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/headers", "get_header",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := BlockHeader(ctx, cli, input.BlockId(), new(beacon.SignedBeaconBlockHeader))
			return err
		})
}

func TestBlockHeaders(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/headers", "get_headers",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := BlockHeaders(ctx, cli, input.Slot, input.ParentRoot, new([]beacon.SignedBeaconBlockHeader))
			return err
		})
}
