package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/shared_test"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
	"testing"
)

func TestPoolAttestations(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/pool", "get_pool_attestations",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			return PoolAttestations(ctx, cli, input.Slot, input.CommitteeIndex, new([]phase0.Attestation))
		})
}

func TestPoolAttesterSlashings(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/pool", "get_pool_attester_slashings",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			return PoolAttesterSlashings(ctx, cli, new([]phase0.AttesterSlashing))
		})
}

func TestPoolProposerSlashings(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/pool", "get_pool_proposer_slashings",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			return PoolProposerSlashings(ctx, cli, new([]phase0.ProposerSlashing))
		})
}

func TestPoolVoluntaryExits(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/pool", "get_pool_voluntary_exits",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			return PoolVoluntaryExits(ctx, cli, new([]phase0.SignedVoluntaryExit))
		})
}

// TODO: test the functions that submit data to the pools with POST methods
