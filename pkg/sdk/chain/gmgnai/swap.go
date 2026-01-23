package gmgnai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strconv"
)

type CreateTransactionParams struct {
	TokenIn       string
	TokenOut      string
	TokenInAmount string
	//TokenOutAmount string
	Wallet    string
	Slippage  float64
	GasFee    float64
	IsAntiMev bool
}

type Quote struct {
	InputMint            string `json:"inputMint"`
	InAmount             string `json:"inAmount"`
	InDecimals           int    `json:"inDecimals"`
	OutDecimals          int    `json:"outDecimals"`
	OutputMint           string `json:"outputMint"`
	OutAmount            string `json:"outAmount"`
	OtherAmountThreshold string `json:"otherAmountThreshold"`
	SwapMode             string `json:"swapMode"`
	SlippageBps          string `json:"slippageBps"`
	PlatformFee          string `json:"platformFee"`
	PriceImpactPct       string `json:"priceImpactPct"`
	RoutePlan            []struct {
		SwapInfo struct {
			AmmKey     string `json:"ammKey"`
			Label      string `json:"label"`
			InputMint  string `json:"inputMint"`
			OutputMint string `json:"outputMint"`
			InAmount   string `json:"inAmount"`
			OutAmount  string `json:"outAmount"`
			FeeAmount  string `json:"feeAmount"`
			FeeMint    string `json:"feeMint"`
		} `json:"swapInfo"`
		Percent int `json:"percent"`
	} `json:"routePlan"`
	TimeTaken float64 `json:"timeTaken"`
}

type RawTx struct {
	SwapTransaction           string `json:"swapTransaction"`
	LastValidBlockHeight      int    `json:"lastValidBlockHeight"`
	PrioritizationFeeLamports int    `json:"prioritizationFeeLamports"`
	RecentBlockhash           string `json:"recentBlockhash"`
}

type CreateTransactionResult struct {
	Code    int
	Message string
	Data    struct {
		Quote        Quote       `json:"quote"`
		RawTx        RawTx       `json:"raw_tx"`
		AmountInUsd  string      `json:"amount_in_usd"`
		AmountOutUsd string      `json:"amount_out_usd"`
		JitoOrderId  interface{} `json:"jito_order_id"`
	}
}

// 有手续费
func (t *Client) CreateTransaction(ctx context.Context, params CreateTransactionParams) (*CreateTransactionResult, error) {

	queryParams := map[string]string{
		"token_in_address":  params.TokenIn,
		"token_out_address": params.TokenOut,
		"in_amount":         params.TokenInAmount,
		"from_address":      params.Wallet,
		"slippage":          "10",
		"swap_mode":         "ExactIn", // ExactIn ExactOut
		"fee":               fmt.Sprintf("%f", params.GasFee),
		"is_anti_mev":       strconv.FormatBool(params.IsAntiMev),
	}

	rsp, err := resty.New().R().SetContext(ctx).
		SetQueryParams(queryParams).
		Get("https://gmgn.ai/defi/router/v1/sol/tx/get_swap_route")
	if err != nil {
		return nil, err
	}

	var result CreateTransactionResult

	err = json.Unmarshal(rsp.Body(), &result)
	if err != nil {
		return nil, err
	}

	return &result, nil

}
