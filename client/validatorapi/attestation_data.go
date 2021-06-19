package validatorapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/zrnt/eth2/beacon/phase0"
)

// Requests that the beacon node produce an AttestationData.
func AttestationData(ctx context.Context, cli eth2api.Client,
	slot common.Slot, committeeIndex common.CommitteeIndex,
	dest *phase0.AttestationData) error {
	q := eth2api.Query{
		"slot":            slot,
		"committee_index": committeeIndex,
	}
	return eth2api.MinimalRequest(ctx, cli, eth2api.QueryGET(q, "/eth/v1/validator/attestation_data"), eth2api.Wrap(dest))
}
