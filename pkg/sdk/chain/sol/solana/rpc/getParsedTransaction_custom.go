package rpc

import (
	"context"
	"store/pkg/sdk/chain/sol/solana"
)

func (t *Client) GetTransactionByID(ctx context.Context, id string) (*GetParsedTransactionResult, error) {

	transaction, err := t.GetParsedTransaction(ctx, solana.MustSignatureFromBase58(id), &GetParsedTransactionOpts{
		Commitment: CommitmentFinalized, MaxSupportedTransactionVersion: &solana.MaxSupportedTransactionVersion,
	})
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
