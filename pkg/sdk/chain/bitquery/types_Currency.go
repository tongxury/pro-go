package bitquery

type Currency struct {
	CollectionAddress    string `json:"CollectionAddress,omitempty"`
	EditionNonce         int    `json:"EditionNonce,omitempty"`
	Fungible             bool   `json:"Fungible,omitempty"`
	IsMutable            bool   `json:"IsMutable,omitempty"`
	Key                  string `json:"Field,omitempty"`
	MetadataAddress      string `json:"MetadataAddress,omitempty"`
	MintAddress          string `json:"MintAddress,omitempty"`
	Name                 string `json:"Name,omitempty"`
	Decimals             int    `json:"Decimals,omitempty"`
	Native               bool   `json:"Native,omitempty"`
	PrimarySaleHappened  bool   `json:"PrimarySaleHappened,omitempty"`
	ProgramAddress       string `json:"ProgramAddress,omitempty"`
	SellerFeeBasisPoints int    `json:"SellerFeeBasisPoints,omitempty"`
	Symbol               string `json:"Symbol,omitempty"`
	TokenCreator         struct {
		Address  []interface{} `json:"Address,omitempty"`
		Share    []interface{} `json:"Share,omitempty"`
		Verified []interface{} `json:"Verified,omitempty"`
	} `json:"TokenCreator,omitempty"`
	TokenStandard      string `json:"TokenStandard,omitempty"`
	UpdateAuthority    string `json:"UpdateAuthority,omitempty"`
	Uri                string `json:"Uri,omitempty"`
	VerifiedCollection bool   `json:"VerifiedCollection,omitempty"`
	Wrapped            bool   `json:"Wrapped,omitempty"`
}

type Currencies []Currency
