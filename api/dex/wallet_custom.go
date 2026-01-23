package dexpb

import (
	chainutils "store/pkg/sdk/chain/utils"
	"store/pkg/sdk/cryptoz"
	helpers "store/pkg/sdk/helper"
)

type EPK string

func (t EPK) PK(secret string) (string, error) {
	decryptString, err := cryptoz.RSA().DecryptString(string(t), secret)
	if err != nil {
		return "", err
	}

	return decryptString, nil
}

func (t EPK) MPK(secret string) string {
	pk, _ := t.PK(secret)
	return pk
}

type Profits []*Profit

func (ts Profits) AsMap() map[string]*Profit {

	mp := make(map[string]*Profit, len(ts))

	for i := range ts {
		x := ts[i]
		mp[x.Token.XId] = x
	}

	return mp
}

// Wallets
type Wallets []*Wallet

func (ts Wallets) AsMapByUserId() map[string]*Wallet {
	result := make(map[string]*Wallet, len(ts))
	for _, t := range ts {
		result[t.UserId] = t
	}

	return result
}

func (ts Wallets) Ids() []string {
	var walletIds []string
	for _, t := range ts {
		walletIds = append(walletIds, t.XId)
	}

	return walletIds
}

func (t *Wallet) WalletName() string {
	return helpers.OrString(t.Name, chainutils.ShortenAddress(t.XId)) + " " +
		helpers.Select(t.IsDefault(), "📌", "")
}

func (t *Wallet) PrivateKey(secret string) string {
	if t == nil {
		return ""
	}

	decryptString, _ := cryptoz.RSA().DecryptString(t.Epk, secret)

	return decryptString
}

func (t *Wallet) PubKey() string {

	return t.XId
}

func (t *Wallet) IsDefault() bool {
	return t.Default
}

// WalletAssets
type WalletPositions []*WalletPosition

func (ts WalletPositions) TokenIds() []string {

	var tokens []string

	for _, x := range ts {
		tokens = append(tokens, x.Token.XId)
	}

	return tokens
}
