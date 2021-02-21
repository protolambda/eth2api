package validatorapi

import (
	"context"
	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon"
)

// Verifies given aggregate and proofs and publishes them on appropriate gossipsub topic.
func PublishAggregateAndProofs(ctx context.Context, cli eth2api.Client,
	aggAndProofs []*beacon.SignedAggregateAndProof) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.BodyPOST("eth/v1/validator/aggregate_and_proofs", aggAndProofs), nil)
}
