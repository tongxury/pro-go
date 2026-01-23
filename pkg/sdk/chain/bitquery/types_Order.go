package bitquery

type Order struct {
	Account     string `json:"Account,omitempty"`
	BuySide     bool   `json:"BuySide,omitempty"`
	LimitAmount string `json:"LimitAmount,omitempty"`
	LimitPrice  string `json:"LimitPrice,omitempty"`
	Mint        string `json:"Mint,omitempty"`
	OrderId     string `json:"OrderId,omitempty"`
	Owner       string `json:"Owner,omitempty"`
	Payer       string `json:"Payer,omitempty"`
}
