package bitquery

type TradeMetadata struct {
	Account     Account  `json:"Account,omitempty"`
	Amount      string   `json:"Amount,omitempty"`
	AmountInUSD string   `json:"AmountInUSD,omitempty"`
	Price       float64  `json:"Price,omitempty"`
	PriceInUSD  float64  `json:"PriceInUSD,omitempty"`
	Currency    Currency `json:"Currency,omitempty"`
}
