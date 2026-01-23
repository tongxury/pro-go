package trader

//
//import (
//	"context"
//	"encoding/json"
//	"errors"
//	"github.com/go-resty/resty/v2"
//	"store/pkg/sdk/chain/sol/solana/rpc"
//)
//
//type PumpFunTrader struct {
//	solana *rpc.Client
//}
//
////func (t PumpFunTrader) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error) {
////
////	if params.Side == "" {
////		params.Side = "buy"
////	}
////
////	if params.Slippage == 0 {
////		params.Slippage = 0.2
////	}
////
////	response, err := resty.New().R().
////		SetContext(ctx).
////		SetBody(map[string]interface{}{
////			"trade_type":    params.Side,
////			"mint":          params.Token,
////			"amount":        params.Amount.Int64(),
////			"slippage":      params.Slippage * 100,
////			"priorityFee":   params.PriorityFee,
////			"userPublicKey": params.Wallet,
////		}).
////		Post("https://pumpapi.fun/api/trade_transaction")
////	if err != nil {
////		return nil, err
////	}
////
////	var result struct {
////		Transaction string
////	}
////
////	if response.IsError() {
////		return nil, errors.New(response.String())
////	}
////
////	if err := json.Unmarshal(response.Body(), &result); err != nil {
////		return nil, err
////	}
////
////	if result.Transaction == "" {
////		return nil, errors.New("response is empty")
////	}
////
////	// 发送交易
////	tx := solana.MustTransactionFromBase58(result.Transaction)
////
////	sign, err := t.solana.SendTransactionAndWaitStatus(ctx,
////		rpc.SendTransactionParams{
////			Transaction: tx,
////			PrivateKeys: []string{params.PrivateKey},
////		})
////	if err != nil {
////		return nil, fmt.Errorf("send transaction err: %v", err)
////	}
////
////	quote, err := t.getQuote(ctx, params)
////	if err != nil {
////		return nil, err
////	}
////
////	return &CreateTransactionResult{
////		TxHash: sign,
////		Quote:  quote,
////	}, nil
////}
//
//func (t PumpFunTrader) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error) {
//
//	if params.Side == "" {
//		params.Side = "buy"
//	}
//
//	if params.Slippage == 0 {
//		params.Slippage = 0.2
//	}
//
//	response, err := resty.New().R().
//		SetContext(ctx).
//		SetBody(map[string]interface{}{
//			"trade_type": params.Side,
//			"mint":       params.Token,
//			"amount":     params.Amount,
//			"slippage":   params.Slippage * 100,
//			//"priorityFee":    params.Priority.Fee,
//			"userPrivateKey": params.PrivateKey,
//		}).
//		Post("https://pumpapi.fun/api/trade")
//	if err != nil {
//		return nil, err
//	}
//
//	var result struct {
//		TxHash string `json:"tx_hash"`
//	}
//
//	if response.IsError() {
//		return nil, errors.New(response.String())
//	}
//
//	if err := json.Unmarshal(response.Body(), &result); err != nil {
//		return nil, err
//	}
//
//	if result.TxHash == "" {
//		return nil, errors.New("response is empty")
//	}
//
//	quote, err := t.getQuote(ctx, params)
//	if err != nil {
//		return nil, err
//	}
//
//	return &CreateTransactionResult{
//		TxHash: result.TxHash,
//		Quote:  quote,
//	}, nil
//}
//
//func (t PumpFunTrader) getQuote(ctx context.Context, params *CreateTransactionParams) (*Quote, error) {
//
//	//price, err := gmgnai.NewClient().GetPrice(ctx, params.Token)
//	//if err != nil {
//	//	return nil, err
//	//}
//	//
//	if params.Side == "buy" {
//
//		return &Quote{
//			Token:          params.Token,
//			SolAmountRaw:   params.Amount.Raw(),
//			TokenAmountRaw: "", // todo
//			Side:           params.Side,
//		}, nil
//
//	}
//
//	return &Quote{
//		Token:          params.Token,
//		TokenAmountRaw: params.Amount.Raw(),
//		SolAmountRaw:   "",
//		Side:           params.Side,
//	}, nil
//
//}
//
//func NewPumpFunTrader(solana *rpc.Client) ITrader {
//	return &PumpFunTrader{solana: solana}
//
//}
