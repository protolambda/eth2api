package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
)

// Retrieves the committees for the given state.
//
// Optional query parameters:
// - epoch: Fetch committees for the given epoch.  If not present then the committees for the epoch of the state will be obtained.
// - index: Restrict returned values to those matching the supplied committee index.
// - slot: Restrict returned values to those matching the supplied slot.
func EpochCommittees(ctx context.Context, cli eth2api.Client,
	stateId eth2api.StateId,
	epoch *beacon.Epoch,
	index *beacon.CommitteeIndex,
	slot *beacon.Slot,
	dest *[]eth2api.Committee) (exists bool, err error) {
	var q eth2api.Query
	if epoch != nil || index != nil || slot != nil {
		q = make(eth2api.Query)
		if epoch != nil {
			q["epoch"] = *epoch
		}
		if index != nil {
			q["index"] = *index
		}
		if slot != nil {
			q["slot"] = *slot
		}
	}
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtQueryGET(q, "eth/v1/beacon/states/%s/committees", stateId.StateId()), eth2api.Wrap(dest))
}

// Returns finality checkpoints for state with given 'stateId'.
// In case finality is not yet achieved, checkpoint should return epoch 0 and ZERO_HASH as root.
func FinalityCheckpoints(ctx context.Context, cli eth2api.Client,
	stateId eth2api.StateId, dest *eth2api.FinalityCheckpoints) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/states/%s/finality_checkpoints", stateId.StateId()), eth2api.Wrap(dest))
}

// Returns Fork object for state with given 'stateId'
func Fork(ctx context.Context, cli eth2api.Client,
	stateId eth2api.StateId, dest *beacon.Fork) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/states/%s/fork", stateId.StateId()), eth2api.Wrap(dest))
}

// Calculates HashTreeRoot for state with given 'stateId'. If stateId is root, same value will be returned.
func StateRoot(ctx context.Context, cli eth2api.Client,
	stateId eth2api.StateId) (root beacon.Root, exists bool, err error) {
	var dest eth2api.RootResponse
	exists, err = eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/states/%s/root", stateId.StateId()), eth2api.Wrap(&dest))
	root = dest.Root
	return
}

// Returns validator specified by state and id or public key along with status and balance.
func StateValidator(ctx context.Context, cli eth2api.Client,
	stateId eth2api.StateId, validatorId eth2api.ValidatorId, dest *eth2api.ValidatorResponse) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("eth/v1/beacon/states/%s/validators/%s", stateId.StateId(), validatorId.ValidatorId()), eth2api.Wrap(dest))
}

// Returns filterable list of validator balances.
//
// Note that any invalid validators with invalid IDs may be ignored, and omitted from the otherwise valid response.
func StateValidatorBalances(ctx context.Context, cli eth2api.Client,
	stateId eth2api.StateId, validatorIds []eth2api.ValidatorId, dest *[]eth2api.ValidatorBalanceResponse) (exists bool, err error) {
	var q eth2api.Query
	if validatorIds != nil {
		q = eth2api.Query{"id": eth2api.ValidatorIdFilter(validatorIds)}
	}
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtQueryGET(q, "eth/v1/beacon/states/%s/validator_balances", stateId.StateId()), eth2api.Wrap(dest))
}

// Returns filterable list of validators with their balance, status and index.
// The status filter is optional, and filters the query to just the given set of status enum values.
//
// Note that any invalid validators with invalid IDs may be ignored, and omitted from the otherwise valid response.
func StateValidators(ctx context.Context, cli eth2api.Client,
	stateId eth2api.StateId, validatorIds []eth2api.ValidatorId, statusFilter []eth2api.ValidatorStatus, dest *[]eth2api.ValidatorResponse) (exists bool, err error) {
	var q eth2api.Query
	if validatorIds != nil || statusFilter != nil {
		q = make(eth2api.Query)
		if validatorIds != nil {
			q["id"] = eth2api.ValidatorIdFilter(validatorIds)
		}
		if statusFilter != nil {
			q["status"] = eth2api.StatusFilter(statusFilter)
		}
	}
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtQueryGET(q, "eth/v1/beacon/states/%s/validators", stateId.StateId()), eth2api.Wrap(dest))
}
