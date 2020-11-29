package eth2api

import "context"

type Beacon struct {
	Client Client
}

func (b *Beacon) Genesis(ctx context.Context, dest *GenesisResponse) (err error) {
	resp := b.Client.Request(ctx, PlainRequest("eth/v1/genesis"))
	if err := resp.Err(); err != nil {
		return err
	}
	err = resp.Decode(dest)
	return
}
