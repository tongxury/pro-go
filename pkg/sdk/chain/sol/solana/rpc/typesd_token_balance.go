package rpc

import "store/pkg/sdk/conv"

type TokenBalance_ struct {
	IsNative    bool          `json:"isNative"`
	Mint        string        `json:"mint"`
	Owner       string        `json:"owner"`
	State       string        `json:"state"`
	TokenAmount UiTokenAmount `json:"tokenAmount"`
}

func (t *TokenBalance_) Val() float64 {
	if t == nil {
		return 0
	}

	return t.TokenAmount.Val()
}

func (t *TokenBalance_) DecimalVal() int64 {
	if t == nil {
		return 0
	}

	return conv.Int64(t.TokenAmount.Amount)
}

func (t *TokenBalance_) Decimals() int64 {

	if t == nil {
		return 0
	}

	return int64(t.TokenAmount.Decimals)
}

type TokenBalances []*TokenBalance_

func (ts TokenBalances) TokenMints() []string {

	var mints []string
	for _, t := range ts {
		mints = append(mints, t.Mint)
	}
	return mints
}

func (ts TokenBalances) TokenDecimals() map[string]int {

	decimals := make(map[string]int, len(ts))
	for _, t := range ts {
		decimals[t.Mint] = int(t.TokenAmount.Decimals)
	}

	return decimals
}
