package bitquery

import (
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
)

type Trade struct {
	Currency       Currency      `json:"Currency,omitempty"`
	Buy            TradeMetadata `json:"Buy,omitempty"`
	Sell           TradeMetadata `json:"Sell,omitempty"`
	Dex            Dex           `json:"Dex,omitempty"`
	Index          int           `json:"Index,omitempty"`
	Market         Market        `json:"Market,omitempty"`
	PriceAsymmetry int           `json:"PriceAsymmetry,omitempty"`
}

// ParsedTrade bitquery 的amount是string类型
type ParsedTrade struct {
	//Account Account

	Token            Currency
	TokenAmount      float64
	TokenAmountInUSD float64
	TokenPrice       float64
	TokenPriceInUSD  float64

	Quote            Currency
	QuoteAmount      float64
	QuoteAmountInUSD float64
	QuotePrice       float64 // 这个是相对于token的价格 基本没用
	QuotePriceInUSD  float64

	IsBuy bool

	Success bool
}

func (t Trade) Cat() (string, string) {

	isStable := func(src TradeMetadata) string {

		if src.Currency.Native {
			return "sol"
		}

		if helper.Contains([]string{"USDC", "USDT"}, src.Currency.Symbol) {
			return "usd"
		}

		if helper.Contains([]string{
			solana.USDTMintString,
			solana.USDCMintString}, src.Currency.MintAddress) {

			return "usd"
		}

		if src.Currency.MintAddress == solana.SolMintString {
			return "sol"
		}

		return ""
	}

	isSellStable := isStable(t.Sell)
	isBuyStable := isStable(t.Buy)

	if isSellStable != "" && isBuyStable != "" {
		return "stable", ""
	}

	if isSellStable != "" {
		return "buy", isSellStable
	}

	if isBuyStable != "" {
		return "sell", isBuyStable
	}
	//
	//// sell
	//if t.Sell.Currency.Native {
	//	return "buy", "sol"
	//}
	//
	//if t.Sell.Currency.MintAddress == solana.SolMintString {
	//	return "buy", "sol"
	//}
	//
	//if t.Sell.Currency.Symbol == "USDC" {
	//	return "buy", "usd"
	//}
	//
	//if helper.SliceContainsAny([]string{t.Sell.Currency.Symbol, "USDC", "USDT"}) {
	//}
	//
	//// buy
	//if t.Buy.Currency.Native {
	//	return "sell", "sol"
	//}
	//
	//if t.Buy.Currency.MintAddress == solana.SolMintString {
	//	return "sell", "sol"
	//}
	//
	//if t.Buy.Currency.Symbol == "USDC" {
	//	return "sell", "usd"
	//}

	return "unstable", ""
}

// Sell 的PriceInUSD 没有值
// Buy 的AmountInUSD 没有值
func (t Trade) Parse() *ParsedTrade {

	cat, quote := t.Cat()

	if cat == "stable" {
		return nil
	}

	if cat == "unstable" {
		return nil
	}

	pt := &ParsedTrade{}

	if cat == "buy" {

		if quote == "sol" {
			pt.Token = t.Buy.Currency
			pt.TokenPrice = t.Buy.Price
			pt.TokenPriceInUSD = t.Buy.PriceInUSD
			pt.TokenAmount = conv.Float64(t.Buy.Amount)
			pt.TokenAmountInUSD = helper.OrNumber(conv.Float64(t.Buy.AmountInUSD), conv.Float64(t.Sell.AmountInUSD))

			pt.Quote = t.Sell.Currency
			pt.QuotePrice = t.Sell.Price
			pt.QuotePriceInUSD = t.Sell.PriceInUSD
			pt.QuoteAmount = conv.Float64(t.Sell.Amount)
			pt.QuoteAmountInUSD = helper.OrNumber(conv.Float64(t.Sell.AmountInUSD), conv.Float64(t.Buy.AmountInUSD))

		}

		if quote == "usd" {
			pt.Token = t.Buy.Currency

			pt.TokenPrice = t.Buy.Price
			pt.TokenPriceInUSD = t.Buy.PriceInUSD
			pt.TokenAmount = conv.Float64(t.Buy.Amount)
			pt.TokenAmountInUSD = helper.OrNumber(conv.Float64(t.Buy.AmountInUSD), conv.Float64(t.Sell.AmountInUSD))

			pt.Quote = t.Sell.Currency
			pt.QuotePrice = t.Sell.Price
			pt.QuotePriceInUSD = 1
			pt.QuoteAmount = conv.Float64(t.Sell.Amount)
			pt.QuoteAmountInUSD = conv.Float64(t.Sell.Amount)

		}

		pt.IsBuy = true

	}

	if cat == "sell" {

		if quote == "sol" {
			pt.Token = t.Sell.Currency
			pt.TokenPrice = t.Sell.Price
			pt.TokenPriceInUSD = t.Sell.PriceInUSD
			pt.TokenAmount = conv.Float64(t.Sell.Amount)
			pt.TokenAmountInUSD = helper.OrNumber(conv.Float64(t.Sell.AmountInUSD), conv.Float64(t.Buy.AmountInUSD))

			pt.Quote = t.Buy.Currency
			pt.QuotePrice = t.Buy.Price
			pt.QuotePriceInUSD = t.Buy.PriceInUSD
			pt.QuoteAmount = conv.Float64(t.Buy.Amount)
			pt.QuoteAmountInUSD = helper.OrNumber(conv.Float64(t.Buy.AmountInUSD), conv.Float64(t.Sell.AmountInUSD))
		}

		if quote == "usd" {
			pt.Token = t.Sell.Currency
			pt.TokenPrice = t.Sell.Price
			pt.TokenPriceInUSD = t.Sell.PriceInUSD
			pt.TokenAmount = conv.Float64(t.Sell.Amount)
			pt.TokenAmountInUSD = helper.OrNumber(conv.Float64(t.Sell.AmountInUSD), conv.Float64(t.Buy.AmountInUSD))

			pt.Quote = t.Buy.Currency
			pt.QuotePrice = t.Buy.Price
			pt.QuotePriceInUSD = 1
			pt.QuoteAmount = conv.Float64(t.Buy.Amount)
			pt.QuoteAmountInUSD = conv.Float64(t.Buy.Amount)
		}

		pt.IsBuy = false
	}

	if pt.TokenPriceInUSD == 0 && pt.TokenAmount != 0 {
		pt.TokenPriceInUSD = pt.QuotePriceInUSD * pt.QuoteAmount / pt.TokenAmount
	}

	if pt.TokenAmountInUSD == 0 && pt.TokenPriceInUSD != 0 {
		pt.TokenAmountInUSD = pt.QuotePriceInUSD * pt.QuoteAmount / pt.TokenPriceInUSD

	}

	if pt.QuotePriceInUSD == 0 && pt.QuoteAmount != 0 {
		pt.QuotePriceInUSD = pt.TokenPriceInUSD * pt.TokenAmount / pt.QuoteAmount
	}

	if pt.QuoteAmountInUSD == 0 && pt.QuotePriceInUSD != 0 {
		pt.QuoteAmountInUSD = pt.TokenPriceInUSD * pt.TokenAmount / pt.QuotePriceInUSD

	}

	return pt
}
