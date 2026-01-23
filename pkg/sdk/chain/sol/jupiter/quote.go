package jupiter

import (
	"context"
	"fmt"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/restyd"
)

type QuoteRequest struct {
	InputMint  string
	OutputMint string
	Amount     string
	//PlatformFeeBps int // 500 as 5%
	SlippageBps  int // 200 as 2%
	AutoSlippage bool
	//MaxAutoSlippageBps            int
	//AutoSlippageCollisionUsdValue int
	SwapMode string // ExactIn | ExactOut
}

// PlatformFee defines model for PlatformFee.
type PlatformFee struct {
	Amount *string `json:"amount,omitempty"`
	FeeBps *int32  `json:"feeBps,omitempty"`
}

type SwapInfo struct {
	AmmKey     string  `json:"ammKey"`
	FeeAmount  string  `json:"feeAmount"`
	FeeMint    string  `json:"feeMint"`
	InAmount   string  `json:"inAmount"`
	InputMint  string  `json:"inputMint"`
	Label      *string `json:"label,omitempty"`
	OutAmount  string  `json:"outAmount"`
	OutputMint string  `json:"outputMint"`
}

type RoutePlanStep struct {
	Percent  int32    `json:"percent"`
	SwapInfo SwapInfo `json:"swapInfo"`
}

type SwapMode string

type QuoteResponse struct {
	ComputedAutoSlippage *int32          `json:"computedAutoSlippage,omitempty"`
	ContextSlot          *float32        `json:"contextSlot,omitempty"`
	InAmount             string          `json:"inAmount"`
	InputMint            string          `json:"inputMint"`
	OtherAmountThreshold string          `json:"otherAmountThreshold"`
	OutAmount            string          `json:"outAmount"`
	OutputMint           string          `json:"outputMint"`
	PlatformFee          *PlatformFee    `json:"platformFee,omitempty"`
	PriceImpactPct       string          `json:"priceImpactPct"`
	RoutePlan            []RoutePlanStep `json:"routePlan"`
	SlippageBps          int32           `json:"slippageBps"`
	SwapMode             SwapMode        `json:"swapMode"`
	TimeTaken            *float32        `json:"timeTaken,omitempty"`
	Error                string          `json:"error,omitempty"`
	ErrorCode            string          `json:"errorCode,omitempty"`
}

func (t *SwapClient) Quote(ctx context.Context, params *QuoteRequest) (*QuoteResponse, error) {

	url := t.endpoint + "/quote"

	//var quote QuoteResponse

	requestParams := map[string]string{
		"inputMint":  params.InputMint,
		"outputMint": params.OutputMint,
		"amount":     params.Amount,
		"swapMode":   helper.OrString(params.SwapMode, "ExactIn"),
		//"useSharedAccounts": "true",
		//"useTokenLedger":           "true", // 加了这个 会出现交易数据过大的问题
		//"skipUserAccountsRpcCalls": "true",
		//"asLegacyTransaction":      "true",  // 加了这个后 会出现找不到路由的情况
	}

	if params.AutoSlippage {
		requestParams["autoSlippage"] = "true"
		//requestParams["autoSlippageCollisionUsdValue"] = "1000"
	} else {
		requestParams["slippageBps"] = helper.OrString(fmt.Sprintf("%d", params.SlippageBps), "50")
	}

	//if params.PlatformFeeBps > 0 {
	//	requestParams["platformFeeBps"] = fmt.Sprintf("%d", params.PlatformFeeBps)
	//}

	quote, err := restyd.Get[QuoteResponse](ctx, restyd.Params{
		Url:    url,
		Params: requestParams,
	})

	if err != nil {
		return nil, err
	}

	//rsp, err := resty.New().R().
	//	SetQueryParams(requestParams).
	//	SetContext(ctx).
	//	Get(url)
	//
	//if err != nil {
	//	return nil, err
	//}

	//if err := json.Unmarshal(rsp.Body(), &quote); err != nil {
	//	return nil, err
	//}

	//if rsp.StatusCode() != 200 {
	//	return nil, errors.New(rsp.String())
	//}

	if quote.ErrorCode != "" {
		if quote.ErrorCode == "COULD_NOT_FIND_ANY_ROUTE" {
			return nil, CouldNotFindAnyRouteError
		}
		if quote.ErrorCode == "NO_ROUTES_ROUND" {
			return nil, CouldNotFindAnyRouteError
		}
		if quote.ErrorCode == "TOKEN_NOT_TRADABLE" {
			return nil, TokenNotTradableError
		}
	}

	return quote, nil
}
