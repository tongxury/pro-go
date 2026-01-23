package raydium

import (
	"context"
	"store/pkg/sdk/chain/sol/solana"
	"testing"
)

func TestSwap(t *testing.T) {

	c := NewClient()

	ctx := context.Background()

	quote, err := c.Quote(ctx, QuoteParams{
		InputMint:  solana.SolMintString,
		OutputMint: "Ho6wN4ff7RdTdXE1UsCZjrjuFVMHyRFTv1oBdbSECnJS",
		Amount:     10000,
		Slippage:   0.05,
	})
	if err != nil {
		return
	}

	c.Swap(ctx, SwapParams{
		ComputeUnitPriceMicroLamports: 1000000,
		SwapResponse:                  quote,
		TxVersion:                     "V0",
		Wallet:                        "3Cyrb9NABtEMJWv5yJ5pnccJPt3UYb1TmyiACu98NQts",
	})

}
