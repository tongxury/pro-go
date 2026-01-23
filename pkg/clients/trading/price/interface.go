package pricing

import "context"

type IPricing interface {
	GetPrice(ctx context.Context, params *GetPriceParams) (*GetPriceResult, error)
}

type GetPriceParams struct {
	Token         string
	TokenDecimals int
}

type GetPriceResult struct {
	Value float64
}
