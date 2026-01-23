package trader

//
//import (
//	"context"
//	"fmt"
//	"github.com/go-resty/resty/v2"
//	"github.com/mr-tron/base58"
//	"net/http"
//	"store/pkg/sdk/chain/sol/solana"
//	"store/pkg/sdk/chain/sol/solana/rpc"
//	"store/pkg/sdk/helper"
//)
//
//type PumpPortalTrader struct {
//	solana *rpc.Client
//}
//
////type CreateTransactionParams struct {
////	// Your wallet public key
////	Wallet string `json:"publicKey"`
////	// "buy" or "sell"
////	Side string `json:"action"`
////	// The contract address of the token you want to trade (this is the text after the '/' in the pump.fun url for the token.)
////	Token string `json:"mint"`
////	// amount: The amount of SOL or tokens to trade. If selling, amount can be a percentage of tokens in your wallet (ex. amount: "100%")
////	Amount float64 `json:"amount"`
////	// denominatedInSol: "true" if amount is SOL, "false" if amount is tokens
////	DenominatedInSol string `json:"denominatedInSol"`
////	// eg. 10
////	Slippage int `json:"slippage"`
////	//priorityFee: Amount to use as priority fee  0.005
////	PriorityFee float64 `json:"priorityFee"`
////	// pool: (optional) Currently 'pump' and 'raydium' are supported options.Default is 'pump'.
////	Pool string `json:"pool"`
////}
//
//func (t PumpPortalTrader) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error) {
//
//	if params.Slippage == 0 {
//		params.Slippage = 0.1
//	}
//
//	result, err := resty.New().R().
//		SetContext(ctx).
//		SetBody(map[string]interface{}{
//			"publicKey":        params.Wallet,
//			"action":           params.Side,
//			"mint":             params.Token,
//			"amount":           params.Amount,
//			"denominatedInSol": helper.Select(params.Side == "buy", "true", "false"),
//			"slippage":         params.Slippage * 100,
//			//"priorityFee":      params.PriorityFee,
//			"pool": "pump",
//		}).
//		Post("https://pumpportal.fun/api/trade-local")
//
//	if err != nil {
//		return nil, nil
//	}
//
//	if result.StatusCode() != http.StatusOK {
//		return nil, fmt.Errorf(result.Status())
//	}
//
//	tx := solana.MustTransactionFromBase58(base58.Encode(result.Body()))
//
//	sign, ok, err := t.solana.SendTransactionAndWaitUtilConfirmed(ctx,
//		rpc.SendTransactionParams{
//			Transaction: tx,
//			PrivateKeys: []string{params.PrivateKey},
//		})
//	if err != nil {
//		return nil, fmt.Errorf("send transaction err: %v", err)
//	}
//
//	return &CreateTransactionResult{
//		TxHash:    sign,
//		Confirmed: ok,
//		Quote:     nil,
//	}, nil
//}
//
//func NewPumpPortalTrader(solana *rpc.Client) ITrader {
//	return &PumpPortalTrader{solana: solana}
//
//}
