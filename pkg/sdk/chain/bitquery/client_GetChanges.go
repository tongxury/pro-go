package bitquery

import (
	"context"
	"time"
)

type GetChangesParams struct {
	Token           string
	IntervalMinutes int // = 5
	Limit           int
}

func (t *Client) GetChanges(ctx context.Context, params *GetChangesParams) (Changes, error) {

	if params.IntervalMinutes == 0 {
		params.IntervalMinutes = 5
	}

	if params.Limit == 0 {
		params.Limit = 10
	}

	sql := `
query Token($token: String, $intervalMinutes: Int, $limit: Int) {
  Solana(dataset: realtime) {
    DEXTradeByTokens(
      orderBy: { descendingByField: "Block_Timefield" }
      where: {
        Trade: {
          Currency: {
            MintAddress: { is: $token}
          }
          PriceAsymmetry: { lt: 0.1 }
        },
        Transaction: {
          Result: {Success: true}
        }
      }
      limit: { count: $limit }
    ) {
      Block {
        Timefield: Time(interval: { in: minutes, count: $intervalMinutes })
      }

      buyVolume: sum(of: Trade_Amount, if: {Trade: {Side: {Type: {is: buy}}}}),
      sellVolume: sum(of: Trade_Amount, if: {Trade: {Side: {Type: {is: sell}}}})
      
      tradeCount: count
      sellCount: count(if: {Trade: {Side: {Type: {is: sell}}}})
      buyCount: count(if: {Trade: {Side: {Type: {is: buy}}}})

      buyers: count(
        distinct: Transaction_Signer
        if: {Trade: {Side: {Type: {is: buy}}}}
      )
      sellers: count(
        distinct: Transaction_Signer
        if: {Trade: {Side: {Type: {is: sell}}}}
      )
      Trade {
        high: PriceInUSD(maximum: Trade_PriceInUSD)
        low: PriceInUSD(minimum: Trade_PriceInUSD)
        open: Price(minimum: Block_Slot)
        close: Price(maximum: Block_Slot)
      }
    }
  }
}
`
	var result DataSolanaResponse[struct {
		DEXTradeByTokens Changes `json:"DEXTradeByTokens"`
	}]

	err := t.Query(ctx, sql,
		map[string]any{
			"token":           params.Token,
			"intervalMinutes": params.IntervalMinutes,
			"limit":           params.Limit,
		},
		&result,
	)
	if err != nil {
		return nil, err
	}

	return result.Data.Solana.DEXTradeByTokens, nil
}

type Change struct {
	Block struct {
		Time time.Time `json:"Timefield"`
	} `json:"Block"`
	Trade struct {
		Close float64 `json:"close"`
		High  float64 `json:"high"`
		Low   float64 `json:"low"`
		Open  float64 `json:"open"`
	} `json:"Trade"`

	TradeCount string `json:"tradeCount"`
	SellCount  string `json:"sellCount"`
	BuyCount   string `json:"buyCount"`

	BuyVolume  string `json:"buyVolume"`
	SellVolume string `json:"sellVolume"`

	Buyers  string `json:"buyers"`
	Sellers string `json:"sellers"`
}

type Changes []*Change
