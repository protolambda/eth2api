package validatorapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Subscribe to a number of sync committee subnets
//
// Sync committees are not present in phase0, but are required for Altair networks.
//
// Subscribing to sync committee subnets is an action performed by VC to enable network participation in Altair networks,
// and only required if the VC has an active validator in an active sync committee.
//
// Err will be non-nil when syncing.
func PrepareSyncCommitteeSubnet(ctx context.Context, cli eth2api.Client,
	signals []*eth2api.SyncCommitteeSubscribeSignal) (syncing bool, err error) {
	req := eth2api.BodyPOST("/eth/v1/validator/sync_committee_subscriptions", signals)
	resp := cli.Request(ctx, req)
	var code uint
	code, err = resp.Decode(nil)
	syncing = code == 503
	return
}
