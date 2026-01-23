package bitquery

type Transaction struct {
	FeePayer string `json:"FeePayer,omitempty"`
	Fee      string `json:"Fee,omitempty"`
	FeeInUSD string `json:"FeeInUSD,omitempty"`
	Index    int    `json:"Index,omitempty"`
	Result   struct {
		ErrorMessage string `json:"ErrorMessage,omitempty"`
		Success      bool   `json:"Success,omitempty"`
	} `json:"Result,omitempty"`
	Signature string `json:"Signature,omitempty"`
	Signer    string `json:"Signer,omitempty"`
}
