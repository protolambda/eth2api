package configapi

import (
	"context"
	"github.com/protolambda/eth2api"
)

// Retrieve deposit contract address and genesis fork version.
func DepositContract(ctx context.Context, cli eth2api.Client, dest *eth2api.DepositContractResponse) error {
	return eth2api.MinimalRequest(ctx, cli, eth2api.PlainGET("/eth/v1/config/deposit_contract"), eth2api.Wrap(dest))
}
