package dexpb

import helpers "store/pkg/sdk/helper"

type TokenMetadatas []*Token_Metadata

type Trades []*Trade

func (ts Trades) Tokens() Tokens {

	mp := map[string]*Token{}
	for _, t := range ts {
		mp[t.Token.XId] = t.Token
		mp[t.QuoteToken.XId] = t.QuoteToken
	}

	return helpers.MapValues(mp)
}

type TradeStates []*TradeState

type OHLCStates []*OHLCState
