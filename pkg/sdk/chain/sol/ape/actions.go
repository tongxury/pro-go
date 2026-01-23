package ape

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

func (t *Client) ListLatestTrades(ctx context.Context, params ListLatestTradesParams) ([]Trade, error) {

	pools, err := t.GetPoolsByToken(ctx, params.Token, true)
	if err != nil {
		return nil, err
	}

	url := "https://api.ape.pro/api/v1/actions/" + pools[0].Id
	result, err := resty.New().R().
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, fmt.Errorf(result.String())
	}

	var trades LatestTrades
	if err := json.Unmarshal(result.Body(), &trades); err != nil {
		return nil, err
	}

	return trades.Txs, nil
}

type ListLatestTradesParams struct {
	Token string
}

type Trade struct {
	BlockId             string    `json:"blockId"`
	Timestamp           time.Time `json:"timestamp"`
	TxHash              string    `json:"txHash"`
	ActionId            string    `json:"actionId"`
	TraderAddress       string    `json:"traderAddress"`
	OfferAsset          string    `json:"offerAsset"`
	OfferAmount         float64   `json:"offerAmount"`
	OfferAssetUsdPrice  float64   `json:"offerAssetUsdPrice"`
	ReturnAsset         string    `json:"returnAsset"`
	ReturnAmount        float64   `json:"returnAmount"`
	ReturnAssetUsdPrice float64   `json:"returnAssetUsdPrice"`
	UsdVolume           float64   `json:"usdVolume"`
	Type                string    `json:"type"`
}
type LatestTrades struct {
	Txs []Trade `json:"txs"`
}
