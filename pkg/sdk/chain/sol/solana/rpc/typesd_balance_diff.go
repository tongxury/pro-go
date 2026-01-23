package rpc

import "store/pkg/sdk/chain/sol/solana"

type BalanceDiff struct {
	Mint     string
	Owner    string
	Decimals int
	Post     string
	Pre      string
}

type BalanceDiffs []*BalanceDiff

func (ts BalanceDiffs) FilterOwner(owner string) BalanceDiffs {

	var result BalanceDiffs
	for _, t := range ts {

		if t.Owner == owner {
			result = append(result, t)
		}
	}

	return result
}

func (ts BalanceDiffs) FilterTokens() BalanceDiffs {

	var result BalanceDiffs
	for _, t := range ts {
		if t.Mint != solana.SolMint.String() {
			result = append(result, t)
		}
	}

	return result
}

func (t *ParsedTransactionMeta) BalanceDiffs() BalanceDiffs {

	tokenDiffMap := map[string]*BalanceDiff{}

	for _, x := range t.PreTokenBalances {

		mint := x.Mint.String()
		var owner string
		if x.Owner != nil {
			owner = x.Owner.String()
		}

		tokenDiffMap[mint+owner] = &BalanceDiff{
			Mint:     mint,
			Owner:    owner,
			Decimals: int(x.UiTokenAmount.Decimals),
			Pre:      x.UiTokenAmount.Amount,
		}
	}
	for _, x := range t.PostTokenBalances {

		mint := x.Mint.String()
		var owner string
		if x.Owner != nil {
			owner = x.Owner.String()
		}

		old := tokenDiffMap[mint+owner]

		if old != nil {
			tokenDiffMap[mint+owner].Post = x.UiTokenAmount.Amount
		} else {
			tokenDiffMap[mint+owner] = &BalanceDiff{
				Mint:     x.Mint.String(),
				Owner:    x.Owner.String(),
				Decimals: int(x.UiTokenAmount.Decimals),
				Post:     x.UiTokenAmount.Amount,
			}
		}
	}

	var balanceDiffs BalanceDiffs
	for _, diff := range tokenDiffMap {
		balanceDiffs = append(balanceDiffs, diff)
	}

	return balanceDiffs
}
