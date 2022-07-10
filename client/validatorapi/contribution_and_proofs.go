package validatorapi

import (
	"context"

	"github.com/protolambda/eth2api"
	"github.com/protolambda/zrnt/eth2/beacon/altair"
)

// Publish multiple signed sync committee contribution and proofs
//
// Err will be non-nil when syncing.
func PublishContributionAndProofs(ctx context.Context, cli eth2api.Client, contribAndProofs []altair.SignedContributionAndProof) error {
	req := eth2api.BodyPOST("/eth/v1/validator/contribution_and_proofs", contribAndProofs)
	resp := cli.Request(ctx, req)
	_, err := resp.Decode(nil)
	return err
}
