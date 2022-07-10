package validatorapi

import (
	"context"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// Aggregates all attestations matching given attestation data root and slot
func AggregateAttestation(ctx context.Context, cli eth2api.Client,
	attDataRoot common.Root, slot common.Slot, dest *phase0.Attestation) error {
	q := eth2api.Query{
		"attestation_data_root": attDataRoot,
		"slot":                  slot,
	}
	return eth2api.MinimalRequest(ctx, cli, eth2api.QueryGET(q, "/eth/v1/validator/aggregate_attestation"), eth2api.Wrap(dest))
}
