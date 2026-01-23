package airwallex

import (
	"context"
	"fmt"
)

func (t *Client) GetPrice(ctx context.Context, priceId string) (*Price, error) {

	var r Price
	err := t.get(ctx, fmt.Sprintf("/api/v1/prices/%s", priceId), nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (t *Client) ListPrices(ctx context.Context) (*List[Price], error) {

	var r List[Price]

	err := t.get(ctx, "/api/v1/prices", nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type Price struct {
	Active       bool           `json:"active"`
	Currency     string         `json:"currency"`
	Description  string         `json:"description"`
	Id           string         `json:"id"`
	Metadata     map[string]any `json:"metadata"`
	Metered      bool           `json:"metered"`
	Name         string         `json:"name"`
	PricingModel string         `json:"pricing_model"`
	ProductId    string         `json:"product_id"`
	Recurring    Recurring      `json:"recurring"`
	Tiers        []struct {
		UpperBound float64 `json:"upper_bound,omitempty"`
		Amount     float64 `json:"amount"`
	} `json:"tiers"`
	Type      string `json:"type"`
	RequestId string `json:"request_id"`
}
