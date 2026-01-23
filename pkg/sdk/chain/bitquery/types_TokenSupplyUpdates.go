package bitquery

type TokenSupplyUpdate struct {
	TokenSupplyUpdate struct {
		Currency    Currency `json:"currency"`
		PostBalance string   `json:"PostBalance"`
		PreBalance  string   `json:"PreBalance"`
	}
}
type TokenSupplyUpdates []*TokenSupplyUpdate
