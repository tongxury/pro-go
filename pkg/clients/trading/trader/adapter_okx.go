package trader

//
//import (
//	"context"
//	"github.com/davecgh/go-spew/spew"
//	"store/pkg/sdk/chain/okxapi"
//	"store/pkg/sdk/chain/sol/solana"
//	"store/pkg/sdk/chain/sol/solana/rpc"
//)
//
//func NewOKXTrader(solana *rpc.Client, okc *okxapi.Client) ITrader {
//	return &OKXTrader{
//		solana: solana,
//		okc:    okc,
//	}
//}
//
//type OKXTrader struct {
//	solana *rpc.Client
//	okc    *okxapi.Client
//}
//
//func (t *OKXTrader) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error) {
//
//	cli := okxapi.NewClient()
//
//	from, to := solana.SolMintString, params.Token
//	if params.Side == "sell" {
//		from, to = params.Token, solana.SolMintString
//	}
//
//	if params.Slippage == 0 {
//		params.Slippage = 0.05
//	}
//
//	swap, err := cli.Swap(ctx, okxapi.SwapParams{
//		Amount:           params.Amount.Raw(),
//		FromTokenAddress: from,
//		ToTokenAddress:   to,
//		Slippage:         params.Slippage,
//		//AutoSlippage:                   "",
//		//MaxAutoSlippage:                "",
//		UserWalletAddress: params.Wallet,
//		//FromTokenReferrerWalletAddress: "",
//		//ToTokenReferrerWalletAddress:   "",
//		//FeePercent:                     "",
//		//GasLimit:                       "",
//		//GasLevel:                       "",
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	tx := solana.MustTransactionFromBase58(swap.Tx.Data)
//
//	spew.Dump(tx)
//
//	sign, ok, err := t.solana.SendTransactionAndWaitUtilConfirmed(ctx,
//		rpc.SendTransactionParams{
//			Transaction: tx,
//			PrivateKeys: []string{params.PrivateKey},
//		})
//
//	return &CreateTransactionResult{
//		TxHash:    sign,
//		Confirmed: ok,
//		Quote: &Quote{
//			Token:          "",
//			TokenAmountRaw: "",
//			SolAmountRaw:   "",
//			Side:           "",
//		},
//	}, err
//}
