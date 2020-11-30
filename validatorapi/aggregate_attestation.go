package validatorapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
)

// Aggregates all attestations matching given attestation data root and slot
func AggregateAttestation(ctx context.Context, cli eth2api.Client,
	attDataRoot beacon.Root, slot beacon.Slot, dest *beacon.Attestation) error {
	q := eth2api.Query{
		"attestation_data_root": attDataRoot,
		"slot":                  slot,
	}
	return eth2api.MinimalRequest(ctx, cli, eth2api.QueryGET(q, "eth/v1/validator/aggregate_attestation"), dest)
}
