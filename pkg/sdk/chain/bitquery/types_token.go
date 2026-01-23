package bitquery

type Token struct {
	MintAddress string  `json:"MintAddress"`
	Name        string  `json:"Name"`
	Symbol      string  `json:"Symbol"`
	Uri         string  `json:"Uri"`
	Decimals    int     `json:"Decimals"`
	Price       float64 `json:"Price"`
	PriceInUSD  float64 `json:"PriceInUSD"`
}

type Tokens []*Token

type TokenMetadata struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	Image       string `json:"image"`
	ShowName    any    `json:"showName"` // true or "true"
	CreatedOn   string `json:"createdOn"`
	Twitter     string `json:"twitter"`
	Website     string `json:"website"`
	
	Websites []struct {
		Label string `json:"label"`
		Url   string `json:"url"`
	} `json:"websites"`
	Socials []struct {
		Url  string `json:"url"`
		Type string `json:"type"`
	} `json:"socials"`
}

type TrendingToken struct {
	Trade       Trade
	TraderCount string `json:"traderCount"`
}

type TrendingTokens []*TrendingToken

func (ts TrendingTokens) Tokens() []string {
	var tokens []string
	for _, t := range ts {
		tokens = append(tokens, t.Trade.Currency.MintAddress)
	}

	return tokens
}
