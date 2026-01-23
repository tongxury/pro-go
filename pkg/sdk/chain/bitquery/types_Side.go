package bitquery

type Side struct {
	Account     Account  `json:"Account"`
	Amount      string   `json:"Amount"`
	AmountInUSD string   `json:"AmountInUSD"`
	Currency    Currency `json:"Currency"`
	Type        string   `json:"Type"`
}

func (s Side) IsBuy() bool {
	return s.Type == "buy"
}
