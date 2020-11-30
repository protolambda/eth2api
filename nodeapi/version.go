package nodeapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Requests that the beacon node identify information about its implementation in a format similar to a  [HTTP User-Agent](https://tools.ietf.org/html/rfc7231#section-5.5.3) field.
func NodeVersion(ctx context.Context, cli eth2api.Client, dest *eth2api.NodeVersionResponse) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("eth/v1/node/version"), dest)
}
