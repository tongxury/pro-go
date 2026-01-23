package okxapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type SwapParams struct {
	// 币种询价数量
	//(数量需包含精度，如兑换 1.00 USDT 需输入 1000000，兑换 1.00 DAI 需输入 1000000000000000000)，币种精度可通过币种列表取得。
	Amount           string // 包含精度
	FromTokenAddress string
	ToTokenAddress   string
	// 滑点限制，最小值：0，最大值：1。（如：0.005代表你接受这笔交易最大 0.5%滑点，1 就代表你接受这笔交易最大 100%的滑点） min:0 max:1
	Slippage float64
	// 默认为 false。当设置为 true 时，原 slippage 参数（如果有传入）将会被 autoSlippage 覆盖，将基于当前市场数据计算并设定自動滑点。
	AutoSlippage string
	// 当 autoSlippage 设置为 true 时，此值为 API 所返回的 autoSlippage 的最大上限，建议用户采用此值以控制风险。
	MaxAutoSlippage   string
	UserWalletAddress string
	///*
	//	fromToken 分佣地址 (支持 SOL 或 SPL Token 分佣，SOL 分佣使用 wallet 地址，SPL Token 分佣使用 token account)
	//	收取分佣费用的 fromToken 地址。使用 API 时，需要结合 feePercent 设置佣金比例，且单笔交易只能选择 fromToken 分佣或 toToken 分佣。
	//*/
	//ReferrerAddress string
	///*
	//	toToken 分佣地址 (只支持 SPL Token 分佣，使用 token account)
	//	收取分佣费用的 toToken 地址。使用 API 时，需要结合 feePercent 设置佣金比例，且单笔交易只能选择 fromToken 分佣或 toToken 分佣。
	//*/
	//ToTokenReferrerAddress string
	/*
		收取 fromToken 分佣费用的钱包地址。新字段针对 SPL Token 不再需要 token account 入参，指定 Sol 钱包地址即可。
		使用 API 时，需要结合 feePercent 设置佣金比例，且单笔交易只能选择 fromToken 分佣或 toToken 分佣。
	*/
	FromTokenReferrerWalletAddress string
	/*
		收取 toToken 分佣费用的钱包地址。新字段针对 SPL Token 不再需要 token account 入参。
		使用 API 时，需要结合 feePercent 设置佣金比例，且单笔交易只能选择 fromToken 分佣或 toToken 分佣。
	*/
	ToTokenReferrerWalletAddress string
	/*
			发送到分佣地址的询价币种数量百分比，最小百分比：0，最大百分比：3。最多支持小数点后 2 位，系统将自动忽略超出的部分。
		  (例如：实际传入 1.326%，但分拥计算时仅会取 1.32% 的分拥比例)
	*/
	FeePercent string
	//SwapReceiverAddress string
	GasLimit string
	// gas价格等级 (默认为 average,交易消耗gas价格水平，可设置为 average、fast 或 slow)
	GasLevel string
}

func (t *Client) ListSupportedChains(ctx context.Context) (any, error) {

	path := "/api/v5/dex/aggregator/supported/chain?chainId=501"

	headers, err := t.authHeaders("GET", path, "")
	if err != nil {
		return 0, err
	}

	result, err := resty.New().R().SetContext(ctx).
		SetHeaders(headers).Get(t.config.Endpoint + path)

	if err != nil {
		return 0, err
	}

	return result.String(), nil
}

func (t *Client) Swap(ctx context.Context, params SwapParams) (*SwapResult, error) {

	//url := t.config.Endpoint
	path := "/api/v5/dex/aggregator/swap" +
		"?chainId=" + CHAIN_ID_SOLANA +
		"&amount=" + params.Amount +
		"&fromTokenAddress=" + params.FromTokenAddress +
		"&toTokenAddress=" + params.ToTokenAddress +
		"&slippage=" + fmt.Sprintf("%f", params.Slippage) +
		"&userWalletAddress=" + params.UserWalletAddress

	headers, err := t.authHeaders("GET", path, "")
	if err != nil {
		return nil, err
	}

	result, err := resty.New().R().SetContext(ctx).
		SetHeaders(headers).Get(t.config.Endpoint + path)

	if err != nil {
		return nil, err
	}

	var response Response[[]SwapResult]

	err = json.Unmarshal(result.Body(), &response)

	if err != nil {
		return nil, err
	}

	if response.Code != "0" {

		return nil, errors.New(result.String())
	}

	if len(response.Data) == 0 {
		return nil, errors.New(result.String())
	}

	return &response.Data[0], nil
}

type SwapResult struct {
	RouterResult struct {
		ChainId       string `json:"chainId"`
		DexRouterList []struct {
			Router        string `json:"router"`
			RouterPercent string `json:"routerPercent"`
			SubRouterList []struct {
				DexProtocol []struct {
					DexName string `json:"dexName"`
					Percent string `json:"percent"`
				} `json:"dexProtocol"`
				FromToken struct {
					Decimal              string `json:"decimal"`
					IsHoneyPot           bool   `json:"isHoneyPot"`
					TaxRate              string `json:"taxRate"`
					TokenContractAddress string `json:"tokenContractAddress"`
					TokenSymbol          string `json:"tokenSymbol"`
					TokenUnitPrice       string `json:"tokenUnitPrice"`
				} `json:"fromToken"`
				ToToken struct {
					Decimal              string `json:"decimal"`
					IsHoneyPot           bool   `json:"isHoneyPot"`
					TaxRate              string `json:"taxRate"`
					TokenContractAddress string `json:"tokenContractAddress"`
					TokenSymbol          string `json:"tokenSymbol"`
					TokenUnitPrice       string `json:"tokenUnitPrice"`
				} `json:"toToken"`
			} `json:"subRouterList"`
		} `json:"dexRouterList"`
		EstimateGasFee string `json:"estimateGasFee"`
		FromToken      struct {
			Decimal              string `json:"decimal"`
			IsHoneyPot           bool   `json:"isHoneyPot"`
			TaxRate              string `json:"taxRate"`
			TokenContractAddress string `json:"tokenContractAddress"`
			TokenSymbol          string `json:"tokenSymbol"`
			TokenUnitPrice       string `json:"tokenUnitPrice"`
		} `json:"fromToken"`
		FromTokenAmount       string `json:"fromTokenAmount"`
		PriceImpactPercentage string `json:"priceImpactPercentage"`
		QuoteCompareList      []struct {
			AmountOut string `json:"amountOut"`
			DexLogo   string `json:"dexLogo"`
			DexName   string `json:"dexName"`
			TradeFee  string `json:"tradeFee"`
		} `json:"quoteCompareList"`
		ToToken struct {
			Decimal              string `json:"decimal"`
			IsHoneyPot           bool   `json:"isHoneyPot"`
			TaxRate              string `json:"taxRate"`
			TokenContractAddress string `json:"tokenContractAddress"`
			TokenSymbol          string `json:"tokenSymbol"`
			TokenUnitPrice       string `json:"tokenUnitPrice"`
		} `json:"toToken"`
		ToTokenAmount string `json:"toTokenAmount"`
		TradeFee      string `json:"tradeFee"`
	} `json:"routerResult"`
	Tx struct {
		Data                 string   `json:"data"`
		From                 string   `json:"from"`
		Gas                  string   `json:"gas"`
		GasPrice             string   `json:"gasPrice"`
		MaxPriorityFeePerGas string   `json:"maxPriorityFeePerGas"`
		MinReceiveAmount     string   `json:"minReceiveAmount"`
		SignatureData        []string `json:"signatureData"`
		Slippage             string   `json:"slippage"`
		To                   string   `json:"to"`
		Value                string   `json:"value"`
	} `json:"tx"`
}
