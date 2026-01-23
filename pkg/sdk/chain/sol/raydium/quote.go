package raydium

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper/mathz"
)

type QuoteParams struct {
	InputMint  string
	OutputMint string
	Amount     int64
	Slippage   float64 // 0.05 => 5%
	TxVersion  string  // V0 Legacy
}

type QuoteResult struct {
	Id      string    `json:"id"`
	Success bool      `json:"success"`
	Version string    `json:"version"`
	Data    Quotation `json:"data"`
}

type Quotation struct {
	SwapType             string  `json:"swapType"`
	InputMint            string  `json:"inputMint"`
	InputAmount          string  `json:"inputAmount"`
	OutputMint           string  `json:"outputMint"`
	OutputAmount         string  `json:"outputAmount"`
	OtherAmountThreshold string  `json:"otherAmountThreshold"`
	SlippageBps          int     `json:"slippageBps"`
	PriceImpactPct       float64 `json:"priceImpactPct"`
	ReferrerAmount       string  `json:"referrerAmount"`
	RoutePlan            []struct {
		PoolId            string        `json:"poolId"`
		InputMint         string        `json:"inputMint"`
		OutputMint        string        `json:"outputMint"`
		FeeMint           string        `json:"feeMint"`
		FeeRate           int           `json:"feeRate"`
		FeeAmount         string        `json:"feeAmount"`
		RemainingAccounts []interface{} `json:"remainingAccounts"`
	} `json:"routePlan"`
}

func (t *Client) Quote(ctx context.Context, params QuoteParams) (*QuoteResult, error) {

	result, err := resty.New().R().SetContext(ctx).SetQueryParams(map[string]string{
		"inputMint":   params.InputMint,
		"outputMint":  params.OutputMint,
		"amount":      conv.Str(params.Amount),
		"slippageBps": conv.Str(mathz.BasicPoints(params.Slippage)),
		"txVersion":   params.TxVersion,
	}).Get(t.endpoint + "/compute/swap-base-in")
	if err != nil {
		return nil, err
	}

	if result.StatusCode() != 200 {
		return nil, errors.New(result.Status())
	}

	var rsp QuoteResult

	err = json.Unmarshal(result.Body(), &rsp)
	if err != nil {
		return nil, err
	}

	return &rsp, nil
}
