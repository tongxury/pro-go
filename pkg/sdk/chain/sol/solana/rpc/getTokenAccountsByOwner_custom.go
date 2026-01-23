package rpc

import (
	"context"
	"store/pkg/sdk/chain/sol/solana"
)

type parsedTokenBalance struct {
	Parsed struct {
		Info struct {
			IsNative    bool           `json:"isNative"`
			Mint        string         `json:"mint"`
			Owner       string         `json:"owner"`
			State       string         `json:"state"`
			TokenAmount *UiTokenAmount `json:"tokenAmount"`
		} `json:"info"`
		Type string `json:"type"`
	} `json:"parsed"`
	Program string `json:"program"`
	Space   int    `json:"space"`
}

type GetTokenBalancesOptions struct {
	Mint     string
	KeepZero bool
}

func (t *Client) GetTokenBalance(ctx context.Context, owner, token string) (*TokenBalance_, error) {

	balances, err := t.GetTokenBalances(ctx, owner, GetTokenBalancesOptions{
		Mint: token,
	})
	if err != nil {
		return nil, err
	}

	if len(balances) == 0 {
		return nil, nil
	}

	return balances[0], nil
}

func (t *Client) GetTokenBalances(ctx context.Context, owner string, options ...GetTokenBalancesOptions) (TokenBalances, error) {

	config := &GetTokenAccountsConfig{
		ProgramId: &solana.TokenProgramID,
	}

	var keepZero bool

	if len(options) > 0 {

		keepZero = options[0].KeepZero

		if options[0].Mint != "" {
			tokenMint := solana.MustPublicKeyFromBase58(options[0].Mint)
			config = &GetTokenAccountsConfig{
				Mint: &tokenMint,
			}
		}
	}

	tokenAccounts, err := t.GetTokenAccountsByOwner(ctx,
		solana.MustPublicKeyFromBase58(owner),
		config,
		&GetTokenAccountsOpts{
			//Commitment: rpc.CommitmentFinalized,
			Encoding: solana.EncodingJSONParsed,
		})

	if err != nil {
		return nil, err
	}

	var result TokenBalances
	for _, x := range tokenAccounts.Value {

		var y struct {
			Parsed struct {
				Info TokenBalance_ `json:"info"`
				Type string        `json:"type"`
			} `json:"parsed"`
			Program string `json:"program"`
			Space   int    `json:"space"`
		}

		err = json.Unmarshal(x.Account.Data.GetRawJSON(), &y)
		if err != nil {
			return nil, err
		}

		if !keepZero {
			if y.Parsed.Info.TokenAmount.Amount != "0" {
				result = append(result, &y.Parsed.Info)
			}
		} else {
			result = append(result, &y.Parsed.Info)
		}

	}

	return result, nil
}
