package rpc

import (
	"store/pkg/sdk/chain/sol/solana"
	"strconv"
)

func (t *UiTokenAmount) Val() float64 {
	if t == nil {
		return 0
	}
	result, _ := strconv.ParseFloat(t.UiAmountString, 64)
	return result
}

type T struct {
	Parsed struct {
		Info struct {
			IsNative    bool          `json:"isNative"`
			Mint        string        `json:"mint"`
			Owner       string        `json:"owner"`
			State       string        `json:"state"`
			TokenAmount UiTokenAmount `json:"tokenAmount"`
		} `json:"info"`
		Type string `json:"type"`
	} `json:"parsed"`
	Program string `json:"program"`
	Space   int    `json:"space"`
}

func (t *ParsedTransaction) MentionsAny(pks ...solana.PK) bool {

	//for i := range tx.Meta.InnerInstructions {
	//	x := tx.Meta.InnerInstructions[i]
	//	for j := range x.Instructions {
	//		y := x.Instructions[j]
	//
	//		y.ProgramId.
	//	}
	//}

	for i := range t.Message.Instructions {

		x := t.Message.Instructions[i]

		if x.ProgramId.IsAnyOf(pks...) {
			return true
		}
	}

	return false

}
