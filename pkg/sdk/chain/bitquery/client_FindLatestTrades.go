package bitquery

import (
	"context"
)

type FindLatestTradesParams struct {
	Token string
	Limit int // = 30
}

func (t *Client) FindLatestTrades(ctx context.Context, params FindLatestTradesParams) (DEXTradeByTokens, error) {

	sql := `
query LatestTrades($token: String, $limit: Int) {
  Solana {
    DEXTradeByTokens(
      orderBy: {descending: Block_Time}
      limit: {count: $limit}
      where: {Trade: {Currency: {MintAddress: {is: $token}}}, Transaction: {Result: {Success: true}}}
    ) {
      Block {
        Time
      }
      Transaction {
        Signature
		Signer
      }
      Trade {
        Market {
          MarketAddress
        }
        Dex {
          ProtocolName
          ProtocolFamily
        }
        AmountInUSD
        Amount
        PriceInUSD
		Price
        Side {
          Type
          Currency {
            Symbol
            MintAddress
            Name
          }
          AmountInUSD
          Amount
        }
      }
    }
  }
}
`

	var result DataSolanaResponse[struct {
		DEXTradeByTokens DEXTradeByTokens `json:"DEXTradeByTokens"`
	}]

	err := t.Query(ctx, sql,
		map[string]any{
			"token": params.Token,
			"limit": 30,
		},
		&result,
	)
	if err != nil {
		return nil, err
	}

	return result.Data.Solana.DEXTradeByTokens, nil
}
