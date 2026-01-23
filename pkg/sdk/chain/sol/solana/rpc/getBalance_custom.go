package rpc

import (
	"context"
	"math"
	"store/pkg/sdk/chain/sol/solana"
)

func (t *Client) GetSolBalance(ctx context.Context, walletPubKey string) (float64, error) {

	balance, err := t.GetBalance(ctx, solana.MustPublicKeyFromBase58(walletPubKey), CommitmentConfirmed)
	if err != nil {
		return 0, err
	}

	return float64(balance.Value) / math.Pow10(9), nil
}
