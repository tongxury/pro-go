package bitquery

type Market struct {
	BaseCurrency  Currency `json:"BaseCurrency,omitempty"`
	MarketAddress string   `json:"MarketAddress,omitempty"`
	QuoteCurrency Currency `json:"QuoteCurrency,omitempty"`
}
