package dexxpb

import "store/pkg/sdk/helper/timed"

func (t *Token) GetMarketCap() float64 {
	return t.GetPrice().GetUsdtValue() * t.GetSupply().GetValue()
}

func (t *Token) Get5mTradeState() *TradeState {

	states := t.GetTradeStates()

	if len(states) == 0 {
		return nil
	}

	if v, ok := states["5m"]; ok {
		return v
	}

	return nil
}

func (t *Token) CreatedFromNow() string {

	if t.CreatedAt == 0 {
		return ""
	}

	return timed.SmartTime(t.CreatedAt)
}

type Tokens []*Token

func (ts Tokens) Prices() Amounts {

	var tokens Amounts

	for _, x := range ts {
		tokens = append(tokens, x.Price)
	}

	return tokens
}

func (ts Tokens) TokenIds() []string {

	var tokens []string

	for _, x := range ts {
		tokens = append(tokens, x.Id)
	}

	return tokens
}

func (ts Tokens) AsMap() map[string]*Token {
	tokens := make(map[string]*Token, len(ts))
	for _, t := range ts {
		tokens[t.Id] = t
	}
	return tokens
}
