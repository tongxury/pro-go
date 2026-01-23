package solscan

type SolscanResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

type Account struct {
	TokenInfo TokenInfo `json:"tokenInfo"`
}

type TokenInfo struct {
	Name     string `json:"name"`
	Symbol   string `json:"symbol"`
	Icon     string `json:"icon"`
	Decimals int    `json:"decimals"`
	//TokenAuthority  interface{} `json:"tokenAuthority"`
	//FreezeAuthority interface{} `json:"freezeAuthority"`
	Supply string `json:"supply"`
	Type   string `json:"type"`
}
