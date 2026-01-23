package jupiter

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
)

type ComputeUnitPriceMicroLamports struct {
	union json.RawMessage
}

func (t *ComputeUnitPriceMicroLamports) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

type PrioritizationFeeLamports struct {
	json.RawMessage
}

//
//func (t *PrioritizationFeeLamports) UnmarshalJSON(b []byte) error {
//	err := t.union.UnmarshalJSON(b)
//	return err
//}

//
//func (t *SwapRequest_PrioritizationFeeLamports) UnmarshalJSON(b []byte) error {
//	err := t.union.UnmarshalJSON(b)
//	return err
//}

/**
userPublicKey (string) — 必填
用户的公钥，表示发起请求的用户。

wrapAndUnwrapSol (boolean) — 默认值为 true
如果为 true，则会自动进行 SOL 的封装和解封装。如果设为 false，则使用 wSOL 代币账户。如果指定了 destinationTokenAccount，此字段将被忽略，因为目标代币账户可能属于不同的用户。

useSharedAccounts (boolean)
启用共享程序账户的使用。这意味着用户无需创建中介代币账户或开放订单账户。如果你使用了 destinationTokenAccount，此字段必须为 true。如果未设置，而路由计划仅通过一个简单的 AMM（非 Openbook 或 Serum），则默认为 false，否则默认为 true。

feeAccount (string)
手续费代币账户，可以是输入铸造或输出铸造。用于计算手续费用，公共密钥由 ["referral_ata", referral_account, mint] 生成，且须确保该账户已创建，不支持 Token2022 代币。

trackingAccount (string)
跟踪账户，可以是任何公钥，便于跟踪交易，特别有利于集成者。可以用特定的 URL 端点查找相关交易。

computeUnitPriceMicroLamports (integer)
优先处理交易的计算单位价格，额外费用将基于 computeUnitLimit（1400000）和该值计算；如果使用 auto，Jupiter 会自动设定优先费用，最高设置为 5,000,000 lamports / 0.005 SOL。

prioritizationFeeLamports (integer)
除签名费用外，支付的优先费用，和 computeUnitPriceMicroLamports 互斥。如果为 auto，Jupiter 会自动设定优先费用，并设上限；另外提供多种设定选项以调整优先费用。

asLegacyTransaction (boolean) — 默认值为 false
请求传统交易，而非默认的版本化交易，需与使用 asLegacyTransaction 的报价配合使用，否则交易可能过大。

useTokenLedger (boolean) — 默认值为 false
在 swap 前的指令中有转账增加输入代币数量时会非常有用，该选项将只使用代币账本余额和后期代币数量之间的差额进行 swap。

destinationTokenAccount (string)
用于接收 swap 后代币的账户公钥。如果未提供，则使用用户的 ATA；如果提供，则假设该账户已初始化。

dynamicComputeUnitLimit (boolean)
如果启用，会进行 swap 模拟，以获取所需的计算单位并设置在 ComputeBudget 的计算单位限制中。启用会稍微增加延迟，因为会有额外的 RPC 调用用于模拟。

skipUserAccountsRpcCalls (boolean)
启用后，将不会对用户账户进行任何 rpc 调用检查。仅当确定所有交易所需账户（如封装或解封装 SOL）已设置好时启用。

dynamicSlippage (object)
基于交易代币类型和用户最大滑点容忍度的数据驱动滑点估算，提供一个最优值以保障用户并确保成功率。

minBps (int32)
用户设定的最小滑点。

maxBps (int32)
用户设定的最大滑点，注意 jup.ag 的用户界面默认值为 300 bps（3%）。

报价响应字段 (quoteResponse)
inputMint (string) — 必填
输入代币的铸造地址。

inAmount (string) — 必填
输入代币的数量。

outputMint (string) — 必填
输出代币的铸造地址。

outAmount (string) — 必填
输出代币的数量。

otherAmountThreshold (string) — 必填
其他金额的阈值，通常用于滑点计算。

swapMode (string) — 必填
可能的值：[ExactIn, ExactOut]，用于指示计算方式。

slippageBps (int32) — 必填
设定的滑点。

platformFee (object)

amount (string)
手续费金额。
feeBps (int32)
手续费的基点。
*/

type SwapRequest struct {
	// AsLegacyTransaction Default is false. Request a legacy transaction rather than the default versioned transaction, needs to be paired with a quote using asLegacyTransaction otherwise the transaction might be too large.
	AsLegacyTransaction *bool `json:"asLegacyTransaction,omitempty"`

	// DestinationTokenAccount Public key of the token account that will be used to receive the token out of the swap. If not provided, the user's ATA will be used. If provided, we assume that the token account is already initialized.
	DestinationTokenAccount *string `json:"destinationTokenAccount,omitempty"`

	// DynamicComputeUnitLimit When enabled, it will do a swap simulation to get the compute unit used and set it in ComputeBudget's compute unit limit. This will increase latency slightly since there will be one extra RPC call to simulate this. Default is `false`.
	DynamicComputeUnitLimit *bool `json:"dynamicComputeUnitLimit,omitempty"`

	// PrioritizationFeeLamports Prioritization fee lamports paid for the transaction in addition to the signatures fee. Mutually exclusive with compute_unit_price_micro_lamports. If `auto` is used, Jupiter will automatically set a priority fee and it will be capped at 5,000,000 lamports / 0.005 SOL.
	PrioritizationFeeLamports *PrioritizationFeeLamports `json:"prioritizationFeeLamports,omitempty"`

	// ProgramAuthorityId The program authority id [0;7], load balanced across the available set by default
	ProgramAuthorityId *int          `json:"programAuthorityId,omitempty"`
	QuoteResponse      QuoteResponse `json:"quoteResponse"`

	// SkipUserAccountsRpcCalls When enabled, it will not do any rpc calls check on user's accounts. Enable it only when you already setup all the accounts needed for the trasaction, like wrapping or unwrapping sol, destination account is already created.
	SkipUserAccountsRpcCalls *bool `json:"skipUserAccountsRpcCalls,omitempty"`

	// UseTokenLedger Default is false. This is useful when the instruction before the swap has a transfer that increases the input token amount. Then, the swap will just use the difference between the token ledger token amount and post token amount.
	UseTokenLedger *bool `json:"useTokenLedger,omitempty"`

	// ComputeUnitPriceMicroLamports The compute unit price to prioritize the transaction, the additional fee will be `computeUnitLimit (1400000) * computeUnitPriceMicroLamports`. If `auto` is used, Jupiter will automatically set a priority fee and it will be capped at 5,000,000 lamports / 0.005 SOL.
	ComputeUnitPriceMicroLamports *ComputeUnitPriceMicroLamports `json:"computeUnitPriceMicroLamports,omitempty"`

	// Tracking account, this can be any public key that you can use to track the transactions, especially useful for integrator. Then, you can use the https://stats.jup.ag/tracking-account/:public-key/YYYY-MM-DD/HH endpoint to get all the swap transactions from this public key.
	TrackingAccount *string `json:"trackingAccount,omitempty"`

	// FeeAccount Fee token account, same as the output token for ExactIn and as the input token for ExactOut, it is derived using the seeds = ["referral_ata", referral_account, mint] and the `REFER4ZgmyYx9c6He5XfaTMiGfdLwRnkV4RPp9t9iF3` referral contract (only pass in if you set a feeBps and make sure that the feeAccount has been created).
	FeeAccount *string `json:"feeAccount,omitempty"`

	// UseSharedAccounts Default is true. This enables the usage of shared program accountns. That means no intermediate token accounts or open orders accounts need to be created for the users. But it also means that the likelihood of hot accounts is higher.
	UseSharedAccounts *bool `json:"useSharedAccounts,omitempty"`

	// UserPublicKey The user public key.
	UserPublicKey string `json:"userPublicKey"`

	// WrapAndUnwrapSol Default is true. If true, will automatically wrap/unwrap SOL. If false, it will use wSOL token account.  Will be ignored if `destinationTokenAccount` is set because the `destinationTokenAccount` may belong to a different user that we have no authority to close.
	WrapAndUnwrapSol *bool `json:"wrapAndUnwrapSol,omitempty"`
}

type SwapResponse struct {
	LastValidBlockHeight      float32  `json:"lastValidBlockHeight"`
	PrioritizationFeeLamports *float32 `json:"prioritizationFeeLamports,omitempty"`
	SwapTransaction           string   `json:"swapTransaction"`
}

func (t *SwapClient) Swap(ctx context.Context, params *SwapRequest) (*SwapResponse, error) {

	url := t.endpoint + "/swap"

	var swap SwapResponse

	rsp, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		Post(url)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(rsp.Body(), &swap)
	if err != nil {
		return nil, err
	}

	if rsp.StatusCode() != 200 {
		return nil, errors.New(rsp.String())
	}

	if swap.SwapTransaction == "" {
		return nil, errors.New("swap transaction is empty")
	}

	//var quote QuoteResponse
	//
	//err = json.Unmarshal(result.Body(), &quote)
	//if err != nil {
	//	return nil, err
	//}

	return &swap, nil
}
