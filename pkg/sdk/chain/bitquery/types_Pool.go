package bitquery

type Pool struct {
	Base struct {
		ChangeAmount string `json:"ChangeAmount"`
		PostAmount   string `json:"PostAmount"`
	} `json:"Base"`
	Dex    Dex    `json:"Dex"`
	Market Market `json:"Market"`
	Quote  struct {
		PostAmount      string  `json:"PostAmount"`
		PostAmountInUSD string  `json:"PostAmountInUSD"`
		PriceInUSD      float64 `json:"PriceInUSD"`
	} `json:"Quote"`
}
