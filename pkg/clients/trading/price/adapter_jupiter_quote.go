package pricing

import (
	"context"
	"fmt"
	"store/pkg/sdk/chain/sol/jupiter"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/helper/mathz/bigdecimal"
)

type JupiterQuote struct {
	c *jupiter.SwapClient
}

func (t JupiterQuote) GetPrice(ctx context.Context, params *GetPriceParams) (*GetPriceResult, error) {

	if params.TokenDecimals == 0 {
		return nil, fmt.Errorf("tokenDecimals must be greater than zero")
	}

	quote, err := t.c.Quote(ctx, &jupiter.QuoteRequest{
		InputMint:  params.Token,
		OutputMint: solana.USDTMint.String(),
		Amount:     bigdecimal.New(1, params.TokenDecimals).Raw(),
	})

	if err != nil {
		return nil, err
	}

	return &GetPriceResult{
		Value: bigdecimal.FromOrigin(quote.OutAmount, 6).Float64(),
	}, nil

}

func NewJupiterQuote(c *jupiter.SwapClient) IPricing {
	return &JupiterQuote{
		c: c,
	}

}
