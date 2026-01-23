package pumpportal

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Client struct {
	endpoint string
}

func New() *Client {
	return &Client{
		endpoint: "https://pumpportal.fun/api/trade-local",
	}
}

// publicKey: Your wallet public key
// action: "buy" or "sell"
// mint: The contract address of the token you want to trade (this is the text after the '/' in the pump.fun url for the token.)
// amount: The amount of SOL or tokens to trade.If selling, amount can be a percentage of tokens in your wallet (ex.amount: "100%")
// denominatedInSol: "true" if amount is SOL, "false" if amount is tokens
// slippage: The percent slippage allowed
// priorityFee: Amount to use as priority fee
// pool: (optional) Currently 'pump' and 'raydium' are supported options.Default is 'pump'.

/*
	response = requests.post(url="https://pumpportal.fun/api/trade-local", data={
	    "publicKey": "Your public key here",
	    "action": "buy",             # "buy" or "sell"
	    "mint": "token CA here",     # contract address of the token you want to trade
	    "amount": 100000,            # amount of SOL or tokens to trade
	    "denominatedInSol": "false", # "true" if amount is amount of SOL, "false" if amount is number of tokens
	    "slippage": 10,              # percent slippage allowed
	    "priorityFee": 0.005,        # amount to use as priority fee
	    "pool": "pump"               # exchange to trade on. "pump" or "raydium"
	})
*/
type CreateTransactionParams struct {
	// Your wallet public key
	Wallet string `json:"publicKey"`
	// "buy" or "sell"
	Side string `json:"action"`
	// The contract address of the token you want to trade (this is the text after the '/' in the pump.fun url for the token.)
	Token string `json:"mint"`
	// amount: The amount of SOL or tokens to trade. If selling, amount can be a percentage of tokens in your wallet (ex. amount: "100%")
	Amount any `json:"amount"`
	// denominatedInSol: "true" if amount is SOL, "false" if amount is tokens
	DenominatedInSol string `json:"denominatedInSol"`
	// eg. 10
	Slippage int `json:"slippage"`
	//priorityFee: Amount to use as priority fee  0.005
	PriorityFee float64 `json:"priorityFee"`
	// pool: (optional) Currently 'pump' and 'raydium' are supported options.Default is 'pump'.
	Pool string `json:"pool"`
}

func (t *Client) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (string, error) {

	if params.Slippage == 0 {
		params.Slippage = 10
	}

	result, err := resty.New().R().
		SetContext(ctx).
		SetBody(params).
		Post(t.endpoint)

	if err != nil {
		return "", nil
	}

	if result.StatusCode() != http.StatusOK {
		return "", fmt.Errorf(result.Status())
	}
	b58 := base64.StdEncoding.EncodeToString(result.Body())

	return b58, nil
}
