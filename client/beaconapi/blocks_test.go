package beaconapi

import (
	"context"
	"testing"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/shared_test"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

func TestBlockAttestations(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/blocks", "get_block_attestations",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := BlockAttestations(ctx, cli, input.BlockId(), new([]phase0.Attestation))
			return err
		})
}

func TestBlock(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/blocks", "get_block",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := Block(ctx, cli, input.BlockId(), new(phase0.SignedBeaconBlock))
			return err
		})
}

// TODO TestBlockV2

func TestBlockRoot(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/blocks", "get_block_root",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, _, err := BlockRoot(ctx, cli, input.BlockId())
			return err
		})
}

func TestPublishBlock(t *testing.T) {
	// TODO: test vectors here don't work yet.
	t.SkipNow()

	shared_test.RunAll(t, "../tests/beacon/blocks", "post_block",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := PublishBlock(ctx, cli, input.Block)
			return err
		})
}
