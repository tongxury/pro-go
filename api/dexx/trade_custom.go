package dexxpb

import mapset "github.com/deckarep/golang-set/v2"

type Trades []*Trade

func (ts Trades) Tokens() Tokens {

	tmp := mapset.NewSet[string]()

	var tokens Tokens
	for _, t := range ts {

		if t.GetToken().GetId() != "" {
			if tmp.Add(t.Token.Id) {
				tokens = append(tokens, t.Token)
			}
		}

		if t.GetQuoteToken().GetId() != "" {
			if tmp.Add(t.QuoteToken.Id) {
				tokens = append(tokens, t.QuoteToken)
			}
		}

	}

	return tokens
}

type TradeSettings_Priorities []*TradeSettings_Priority

func (ts TradeSettings_Priorities) AsMap() map[string]*TradeSettings_Priority {

	priorities := make(map[string]*TradeSettings_Priority, len(ts))

	for _, x := range ts {
		priorities[x.Value] = x
	}

	return priorities
}
