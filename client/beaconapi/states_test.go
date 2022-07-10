package beaconapi

import (
	"context"
	"testing"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/eth2api/shared_test"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

func TestEpochCommittees(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/states", "get_committees",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := EpochCommittees(ctx, cli,
				input.StateId(), input.Epoch, input.CommitteeIndex, input.Slot, new([]eth2api.Committee))
			return err
		})
}

// TODO TestSyncCommittees

func TestFinalityCheckpoints(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/states", "get_finality_checkpoints",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := FinalityCheckpoints(ctx, cli, input.StateId(), new(eth2api.FinalityCheckpoints))
			return err
		})
}

func TestFork(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/states", "get_fork",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := Fork(ctx, cli, input.StateId(), new(common.Fork))
			return err
		})
}

func TestStateRoot(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/states", "get_state_root",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, _, err := StateRoot(ctx, cli, input.StateId())
			return err
		})
}

func TestStateValidator(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/states", "get_validator",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := StateValidator(ctx, cli, input.StateId(), input.ValidatorId(), new(eth2api.ValidatorResponse))
			return err
		})
}

func TestStateValidatorBalances(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/states", "get_validator_balances",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := StateValidatorBalances(ctx, cli, input.StateId(), input.ValidatorIds(), new([]eth2api.ValidatorBalanceResponse))
			return err
		})
}

func TestStateValidators(t *testing.T) {
	shared_test.RunAll(t, "../tests/beacon/states", "get_validators",
		func(ctx context.Context, input *shared_test.Input, cli eth2api.Client) error {
			_, err := StateValidators(ctx, cli, input.StateId(), input.ValidatorIds(), input.StatusFilter, new([]eth2api.ValidatorResponse))
			return err
		})
}
