package bitquery

type DataSolanaResponse[T any] struct {
	Data struct {
		Solana T
	} `json:"data"`
}

type SubscriptionDataSolana[T any] struct {
	Payload struct {
		Data struct {
			Solana T
		} `json:"data"`
	} `json:"payload"`
}

type DEXTradesWrapper = struct {
	DEXTrades DEXTrades `json:"DEXTrades,omitempty"`
}
type DEXPoolsWrapper = struct {
	DEXPools DEXPools `json:"DEXPools,omitempty"`
}

type TokenSupplyUpdatesWrapper = struct {
	TokenSupplyUpdates TokenSupplyUpdates `json:"TokenSupplyUpdates,omitempty"`
}
