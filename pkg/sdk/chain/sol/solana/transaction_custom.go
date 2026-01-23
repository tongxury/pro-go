package solana

import (
	"encoding/base64"
	"fmt"
	"github.com/mr-tron/base58"
)

type AddressLookupTable struct {
	mp map[PublicKey]PublicKeySlice
}

func NewAddressLookupTable() AddressLookupTable {
	return AddressLookupTable{
		mp: make(map[PublicKey]PublicKeySlice),
	}
}

func (t AddressLookupTable) Add(programId string, accounts ...string) {
	for _, x := range accounts {
		t.mp[MPK(programId)] = append(t.mp[MPK(programId)], MPK(x))
	}
}

func (t AddressLookupTable) Value() map[PublicKey]PublicKeySlice {
	return t.mp
}

func MustTransactionFromBase58(b58 string) *Transaction {
	data, err := base58.Decode(b58)
	if err != nil {
		panic(err)
	}
	fromBytes, err := TransactionFromBytes(data)
	if err != nil {
		panic(err)
	}

	return fromBytes
}

func MustTransactionFromBase64(b64 string) *Transaction {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		panic(err)
	}
	fromBytes, err := TransactionFromBytes(data)
	if err != nil {
		panic(err)
	}

	return fromBytes
}

func (tx *Transaction) Signing(privateKey string) (Signature, error) {

	txMessageBytes, err := tx.Message.MarshalBinary()
	if err != nil {
		return zeroSignature, fmt.Errorf("could not serialize transaction: %w", err)
	}

	//txMessageBytes, err := tx.Message.MarshalBinary()
	//if err != nil {
	//	return fmt.Errorf("could not serialize transaction: %w", err)
	//}
	//
	//signature, err := w.PrivateKey.Sign(txMessageBytes)
	//if err != nil {
	//	return fmt.Errorf("could not sign transaction: %w", err)
	//}
	//
	////tx.Signatures 默认是 64个1
	//tx.Signatures = append(tx.Signatures, signature)
	//
	//if len(tx.Signatures) > 1 && tx.Signatures[0].IsZero() {
	//	tx.Signatures = tx.Signatures[1:]
	//}

	pk, err := PrivateKeyFromBase58(privateKey)
	if err != nil {
		return zeroSignature, err
	}

	signature, err := pk.Sign(txMessageBytes)

	if err != nil {
		return zeroSignature, fmt.Errorf("could not sign transaction: %w", err)
	}

	tx.Signatures = []Signature{signature}
	//fmt.Println(signature.String())
	//
	////tx.Signatures 默认是 zeroSignature(64个1)
	//tx.Signatures = append(tx.Signatures, signature)
	//
	//if len(tx.Signatures) > 1 && tx.Signatures[0].IsZero() {
	//	tx.Signatures = tx.Signatures[1:]
	//}

	return signature, nil
}
