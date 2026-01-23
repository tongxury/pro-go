package raydium

import (
	"context"
	"store/pkg/sdk/chain/sol/solana"
	"testing"
)

func TestQuote(t *testing.T) {

	c := NewClient()

	ctx := context.Background()

	c.Quote(ctx, QuoteParams{
		InputMint:  solana.SolMintString,
		OutputMint: "Ho6wN4ff7RdTdXE1UsCZjrjuFVMHyRFTv1oBdbSECnJS",
		Amount:     10000,
		Slippage:   0.05,
	})

}
