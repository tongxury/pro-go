package pumpportal

type TokenCreate struct {
	Signature             string  `json:"signature"`
	Mint                  string  `json:"mint"`
	TraderPublicKey       string  `json:"traderPublicKey"`
	TxType                string  `json:"txType"`
	InitialBuy            float64 `json:"initialBuy"`
	SolAmount             float64 `json:"solAmount"`
	BondingCurveKey       string  `json:"bondingCurveKey"`
	VTokensInBondingCurve float64 `json:"vTokensInBondingCurve"`
	VSolInBondingCurve    float64 `json:"vSolInBondingCurve"`
	MarketCapSol          float64 `json:"marketCapSol"`
	Name                  string  `json:"name"`
	Symbol                string  `json:"symbol"`
	Uri                   string  `json:"uri"`
}
