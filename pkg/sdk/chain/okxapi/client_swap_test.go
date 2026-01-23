package okxapi

import (
	"context"
	"store/pkg/sdk/chain/sol/solana"
	"testing"
)

func TestSwap(t *testing.T) {

	ctx := context.Background()

	c := NewClient()

	c.Swap(ctx, SwapParams{
		Amount:           "100000",
		FromTokenAddress: solana.SolMintString,
		ToTokenAddress:   "mF94n83fqzRDNvs9wQnKAQ62e88fe1JhZwg9MiBpump",
		Slippage:         0.005,
		AutoSlippage:     "true",
		//MaxAutoSlippage:                "",
		UserWalletAddress: "3Cyrb9NABtEMJWv5yJ5pnccJPt3UYb1TmyiACu98NQts",
		//FromTokenReferrerWalletAddress: "",
		//ToTokenReferrerWalletAddress:   "",
		//FeePercent:                     "",
		//GasLimit:                       "",
		//GasLevel:                       "",
	})

	//c.ListSupportedChains(ctx)

}
