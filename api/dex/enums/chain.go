package enums

import (
	tgbotpb "store/api/dex"
	"store/pkg/sdk/chain/sol/solana"
)

var (
	TokenSolana = &tgbotpb.Token{
		XId: solana.SolMint.String(),
		//Symbol:   "SOL",
		//Decimals: 9,
		Metadata: &tgbotpb.Token_Metadata{
			Name:        "Solana",
			Symbol:      "SOL",
			Description: "",
			Image:       "https://raw.githubusercontent.com/solana-labs/token-list/main/assets/mainnet/So11111111111111111111111111111111111111112/logo.png",
			CreatedOn:   "",
			Twitter:     "",
			Website:     "",
			Uri:         "",
			XId:         solana.SolMint.String(),
			Decimals:    9,
		},
	}
)

func NativeTokenById(id string) *tgbotpb.Token {

	switch id {
	case TokenSolana.XId:
		return TokenSolana
	default:
		return TokenSolana
	}
}

var ChainSolana = &tgbotpb.Chain{
	XId: "solana",
	//NativeToken: TokenSolana,
}

func ChainByName(name string) *tgbotpb.Chain {

	return ChainSolana
}

//const (
//	ChainTxStatus_Pending   string = "pending"
//	ChainTxStatus_Confirmed string = "confirmed"
//	ChainTxStatus_Failed    string = "failed"
//	ChainTxStatus_NotFound  string = "notFound"
//)
