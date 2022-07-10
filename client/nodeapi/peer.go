package nodeapi

import (
	"context"

	"github.com/protolambda/eth2api"
)

// Retrieves data about the given peer
func Peer(ctx context.Context, cli eth2api.Client, peerId eth2api.ApiPeerId, dest *eth2api.Peer) (exists bool, err error) {
	return eth2api.SimpleRequest(ctx, cli, eth2api.FmtGET("/eth/v1/node/peers/%s", peerId), eth2api.Wrap(dest))
}

// Retrieves number of known peers.
func PeerCount(ctx context.Context, cli eth2api.Client, dest *eth2api.PeerCountResponse) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("/eth/v1/node/peer_count"), eth2api.Wrap(dest))
}

// Retrieves data about the node's network peers. By default this returns all peers. Different query params are combined using AND conditions. Array items with OR conditions.
func Peers(ctx context.Context, cli eth2api.Client, state []eth2api.ConnectionState, direction []eth2api.ConnectionDirection, dest *[]eth2api.Peer) error {
	var q eth2api.Query
	if state != nil || direction != nil {
		q = make(eth2api.Query)
		if state != nil {
			q["state"] = state
		}
		if direction != nil {
			q["direction"] = direction
		}
	}
	// TODO: current spec has half-removed pagination. A "meta" field alongside the "data" field.
	return eth2api.MinimalRequest(ctx, cli, eth2api.QueryGET(q, "/eth/v1/node/peers"), eth2api.Wrap(dest))
}
