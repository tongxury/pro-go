package bitquery

type Account struct {
	Address string `json:"Address,omitempty"`
	Owner   string `json:"Owner,omitempty"`
	Token   struct {
		Owner string `json:"Owner,omitempty"`
	} `json:"Token,omitempty"`
}
