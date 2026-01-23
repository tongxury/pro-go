package bitquery

import (
	"context"
	"time"
)

func (t *Client) GetTrendingTokens(ctx context.Context, since time.Time, limit int) (TrendingTokens, error) {

	sql := `
query TrendingTokens ($since: DateTime, $limit: Int)  {
  Solana {
    DEXTradeByTokens(
      limit: { count: $limit }
      orderBy: { descendingByField: "traderCount" }
      where: {
         Trade: {Side: {Type: {is: buy}}, Currency: {MintAddress: {not: "So11111111111111111111111111111111111111112"}}},
		 Block: {Time: {since: $since}}
     }
    ) {
      Trade {
        Currency {
          Name
          Symbol
          MintAddress
        }
      }
      traderCount: count(distinct: Transaction_Signer)
    }
  }
}
`

	var result DataSolanaResponse[struct {
		DEXTradeByTokens TrendingTokens `json:"DEXTradeByTokens"`
	}]

	err := t.Query(ctx, sql,
		map[string]any{
			"since": since.UTC().Format("2006-01-02T15:04:05Z"),
			"limit": limit,
		},
		&result,
	)
	if err != nil {
		return nil, err
	}

	return result.Data.Solana.DEXTradeByTokens, nil
}
