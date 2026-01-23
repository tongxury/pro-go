package bitquery

import (
	"context"
	"store/pkg/sdk/conv"
)

func (t *Client) GetTokenSupplies(ctx context.Context, tokens []string) (map[string]float64, error) {

	sql := `
query LatestTrades($tokens: [String!]) {
    Solana {
    TokenSupplyUpdates(
      limit:{count:1}
      orderBy:{descending:Block_Time}
      where: {TokenSupplyUpdate: {Currency: {MintAddress: {in: $tokens}}}}
    ) {
      TokenSupplyUpdate {
        Amount
        Currency {
          MintAddress
          Name
          Decimals
          Symbol
          Uri
        }
        PreBalance
        PostBalance
      }
    }
  }
}
`

	var result DataSolanaResponse[struct {
		TokenSupplyUpdates TokenSupplyUpdates `json:"TokenSupplyUpdates"`
	}]

	err := t.Query(ctx, sql,
		map[string]any{
			"tokens": tokens,
		},
		&result,
	)
	if err != nil {
		return nil, err
	}

	mp := make(map[string]float64, len(tokens))

	for _, x := range result.Data.Solana.TokenSupplyUpdates {
		mp[x.TokenSupplyUpdate.Currency.MintAddress] = conv.Float64(x.TokenSupplyUpdate.PostBalance)
	}

	return mp, nil
}
