package bitquery

type DEXTradeByToken struct {
	Block       Block       `json:"Block"`
	Trade       Trade       `json:"Trade"`
	Transaction Transaction `json:"Transaction"`
}

type DEXTradeByTokens []DEXTradeByToken
