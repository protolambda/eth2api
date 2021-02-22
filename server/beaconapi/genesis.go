package beaconapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Serve details of the chain's genesis which can be used to identify chain.
func Genesis(backend *BeaconBackend) eth2api.Route {
	return eth2api.MakeRoute(eth2api.GET, "eth/v1/beacon/genesis",
		func(ctx context.Context, req eth2api.Request) eth2api.PreparedResponse {
			genesis := backend.Chain.Genesis()
			out := eth2api.GenesisResponse{
				GenesisTime:           genesis.Time,
				GenesisValidatorsRoot: genesis.ValidatorsRoot,
				GenesisForkVersion:    backend.Spec.GENESIS_FORK_VERSION,
			}
			return eth2api.RespondOK(eth2api.Wrap(&out))
		})
}
