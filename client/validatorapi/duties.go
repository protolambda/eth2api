package validatorapi

import (
	"context"
	"fmt"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
)

// Requests the beacon node to provide a set of attestation duties,
// which should be performed by the requested validators, for a particular epoch.
//
// Duties should only need to be checked once per epoch,
// however a chain reorganization (of > MIN_SEED_LOOKAHEAD epochs) could occur,
// resulting in a change of duties. For full safety, you should monitor head events and confirm the
// dependent root in this response matches:
//
//   - event.previous_duty_dependent_root when `compute_epoch_at_slot(event.slot) == epoch`
//
//   - event.current_duty_dependent_root when `compute_epoch_at_slot(event.slot) + 1 == epoch`
//
//   - event.block otherwise
//
// The dependent_root value is `get_block_root_at_slot(state, compute_start_slot_at_epoch(epoch - 1) - 1)`
// or the genesis block root in the case of underflow.
//
// Err will be non-nil when syncing.
func AttesterDuties(ctx context.Context, cli eth2api.Client,
	epoch common.Epoch, indices []common.ValidatorIndex, dest *eth2api.DependentAttesterDuties) (syncing bool, err error) {
	req := eth2api.BodyPOST(fmt.Sprintf("/eth/v1/validator/duties/attester/%d", epoch), indices)
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(dest) // not wrapped, the request type already breaks the `data` boundary
	syncing = code == 503
	return
}

// Request beacon node to provide all validators that are scheduled to propose a block in the given epoch.
//
// Duties should only need to be checked once per epoch,
// however a chain reorganization could occur that results in a change of duties. For full safety,
// you should monitor head events and confirm the dependent root in this response matches:
//
//  - event.current_duty_dependent_root when `compute_epoch_at_slot(event.slot) == epoch`
//
//  - event.block otherwise
//
// The dependent_root value is `get_block_root_at_slot(state, compute_start_slot_at_epoch(epoch) - 1)`
// or the genesis block root in the case of underflow.
//
// Err will be non-nil when syncing.
func ProposerDuties(ctx context.Context, cli eth2api.Client, epoch common.Epoch, dest *eth2api.DependentProposerDuty) (syncing bool, err error) {
	req := eth2api.FmtGET("/eth/v1/validator/duties/proposer/%d", epoch)
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(dest) // not wrapped, the request type already breaks the `data` boundary
	syncing = code == 503
	return
}

// Requests the beacon node to provide a set of sync committee duties for a particular epoch.
func SyncCommitteeDuties(ctx context.Context, cli eth2api.Client,
	epoch common.Epoch, indices []common.ValidatorIndex, dest *[]eth2api.SyncCommitteeDuty) (syncing bool, err error) {
	req := eth2api.BodyPOST(fmt.Sprintf("/eth/v1/validator/duties/sync/%d", epoch), indices)
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(eth2api.Wrap(dest))
	syncing = code == 503
	return
}
