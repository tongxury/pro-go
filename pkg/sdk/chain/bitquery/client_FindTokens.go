package bitquery

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
)

type FindTokensParams struct {
	Ids      []string
	NameLike string
}

type T struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Websites    []struct {
		Label string `json:"label"`
		Url   string `json:"url"`
	} `json:"websites"`
	Socials []struct {
		Url  string `json:"url"`
		Type string `json:"type"`
	} `json:"socials"`
}

func (t *Client) GetTokenMetadata(ctx context.Context, url string) (*TokenMetadata, error) {

	rsp, err := resty.New().R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}

	if !rsp.IsSuccess() {
		return nil, errors.New(rsp.Status())
	}

	if !strings.Contains(rsp.Header().Get("Content-Type"), "application/json") {
		return nil, errors.New(rsp.Status())
	}

	var m TokenMetadata
	err = json.Unmarshal(rsp.Body(), &m)
	if err != nil {
		return nil, err
	}

	return &m, nil

}

func (t *Client) GetTokenById(ctx context.Context, id string) (*Token, error) {

	tokens, err := t.FindTokens(ctx, FindTokensParams{Ids: []string{id}})
	if err != nil {
		return nil, err
	}

	if len(tokens) == 0 {
		return nil, fmt.Errorf("token not found")
	}

	return tokens[0], nil
}

func (t *Client) FindTokens(ctx context.Context, params FindTokensParams) (Tokens, error) {

	sql := `
query Token($tokens: [String!], $nameLike: String) {
  Solana {
    DEXTrades(
      orderBy: {descending: Block_Time}
      limit: {count: 100}
      limitBy: {by: Trade_Buy_Currency_MintAddress, count: 1}
      where: { Trade: {Buy: {Currency: {MintAddress: {in: $tokens}, Name: {includes: $nameLike}}}}}) {
      Trade {
        Buy {
          Price
          PriceInUSD
          Currency {
            Name
            Symbol
            MintAddress
			Uri
			Decimals
          }
        }
        Sell {
          Currency {
            Name
            Symbol
            MintAddress
          }
        }
      }
    }
  }
}
`

	var result DataSolanaResponse[struct {
		DEXTrades DEXTrades `json:"DEXTrades"`
	}]

	err := t.Query(ctx, sql,
		map[string]any{
			"tokens":   params.Ids,
			"nameLike": params.NameLike,
		},
		&result,
	)
	if err != nil {
		return nil, err
	}

	var tokens Tokens
	for _, x := range result.Data.Solana.DEXTrades {
		tokens = append(tokens, &Token{
			MintAddress: x.Trade.Buy.Currency.MintAddress,
			Name:        x.Trade.Buy.Currency.Name,
			Symbol:      x.Trade.Buy.Currency.Symbol,
			Uri:         x.Trade.Buy.Currency.Uri,
			Decimals:    x.Trade.Buy.Currency.Decimals,
			Price:       x.Trade.Buy.Price,
			PriceInUSD:  x.Trade.Buy.PriceInUSD,
		})

	}

	return tokens, nil
}
