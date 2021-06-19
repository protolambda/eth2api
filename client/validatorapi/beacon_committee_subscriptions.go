package validatorapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Signal beacon node to prepare for a committee subnet.
//
// After beacon node receives this request,
// search using discv5 for peers related to this subnet
// and replace current peers with those ones if necessary.
//
// If validator `is_aggregator`, beacon node must:
// - announce subnet topic subscription on gossipsub
// - aggregate attestations received on that subnet
//
// Success when the slot signature is valid and beacon node has prepared the attestation subnet.
// Note that, we cannot be certain Beacon node will find peers for that subnet for various reasons.
//
// Err will be non-nil when syncing.
func PrepareBeaconCommitteeSubnet(ctx context.Context, cli eth2api.Client,
	signals []*eth2api.BeaconCommitteeSubscribeSignal) (syncing bool, err error) {
	req := eth2api.BodyPOST("/eth/v1/validator/beacon_committee_subscriptions", signals)
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(nil)
	syncing = code == 503
	return
}
