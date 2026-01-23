package dexpb

import mapset "github.com/deckarep/golang-set/v2"

type TradeSettings_Priorities []*TradeSettings_Priority

func (ts TradeSettings_Priorities) AsMap() map[string]*TradeSettings_Priority {

	mp := make(map[string]*TradeSettings_Priority)
	for _, x := range ts {
		mp[x.Value] = x
	}

	return mp
}

// Orders
type Orders []*Order

func (ts Orders) Ids() []string {

	var ids []string
	for _, t := range ts {
		ids = append(ids, t.XId)
	}
	return ids
}

func (ts Orders) GroupByToken() map[string]Orders {

	result := make(map[string]Orders, len(ts))
	for _, t := range ts {
		result[t.GetToken().GetXId()] = append(result[t.GetToken().GetXId()], t)
	}
	return result
}

func (ts Orders) Tokens() Tokens {

	tokenIdSet := mapset.NewSet[string]()

	var tokens Tokens

	for _, t := range ts {
		if tokenIdSet.Add(t.Token.GetXId()) {
			tokens = append(tokens, t.Token)
		}
	}

	return tokens
}

func (ts Orders) TokenIds() []string {

	tokenIdSet := mapset.NewSet[string]()

	for _, t := range ts {
		tokenIdSet.Add(t.Token.GetXId())
	}

	return tokenIdSet.ToSlice()
}

func (ts Orders) UserIds() []string {
	var wallets []string
	for _, t := range ts {
		wallets = append(wallets, t.User.XId)
	}

	return mapset.NewSet(wallets...).ToSlice()
}

func (ts Orders) Wallets() []string {
	var wallets []string
	for _, t := range ts {
		wallets = append(wallets, t.Wallet.XId)
	}

	return mapset.NewSet(wallets...).ToSlice()
}
