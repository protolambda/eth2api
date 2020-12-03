package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/shared_test"
	"github.com/protolambda/zrnt/eth2/beacon"
	"testing"
)

func TestPoolAttestations(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/pool", "get_pool_attestations",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			return PoolAttestations(ctx, cli, input.Slot, input.CommitteeIndex, new([]beacon.Attestation))
		})
}

func TestPoolAttesterSlashings(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/pool", "get_pool_attester_slashings",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			return PoolAttesterSlashings(ctx, cli, new([]beacon.AttesterSlashing))
		})
}

func TestPoolProposerSlashings(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/pool", "get_pool_proposer_slashings",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			return PoolProposerSlashings(ctx, cli, new([]beacon.ProposerSlashing))
		})
}

func TestPoolVoluntaryExits(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/pool", "get_pool_voluntary_exits",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			return PoolVoluntaryExits(ctx, cli, new([]beacon.SignedVoluntaryExit))
		})
}

// TODO: test the functions that submit data to the pools with POST methods
