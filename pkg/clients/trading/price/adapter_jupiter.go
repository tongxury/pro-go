package pricing

import (
	"context"
	"store/pkg/sdk/chain/sol/jupiter"
	"store/pkg/sdk/conv"
)

type Jupiter struct {
	c *jupiter.JupiterClient
}

func (t Jupiter) GetPrice(ctx context.Context, params *GetPriceParams) (*GetPriceResult, error) {

	prices, err := t.c.GetPrices(ctx, jupiter.GetPricesParams{
		PublicKeys: []string{params.Token},
	})
	if err != nil {
		return nil, err
	}

	if prices.Data[params.Token] == nil {
		return &GetPriceResult{}, nil
	}

	return &GetPriceResult{
		Value: conv.Float64(prices.Data[params.Token].Price),
	}, nil
}

func NewJupiter(c *jupiter.JupiterClient) IPricing {
	return &Jupiter{
		c: c,
	}

}
