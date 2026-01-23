package bitquery

type DEXTrade struct {
	Block       Block       `json:"Block,omitempty"`
	ChainId     string      `json:"ChainId,omitempty"`
	Trade       Trade       `json:"Trade,omitempty"`
	Transaction Transaction `json:"Transaction,omitempty"`
}

type DEXTrades []*DEXTrade
