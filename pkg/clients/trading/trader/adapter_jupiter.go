package trader

import (
	"context"
	"math"
	"store/pkg/sdk/chain/sol/jupiter"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/chain/sol/solana/rpc"
	"store/pkg/sdk/conv"
	"strings"
)

type JupiterTrader struct {
	solana *rpc.Client
	jc     *jupiter.SwapClient
}

func (t JupiterTrader) CreateInstructions(ctx context.Context, params *CreateCreateInstructionsParams) (*CreateInstructionsResult, error) {
	//TODO implement me
	panic("implement me")
}

func (t JupiterTrader) CreateRawTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateRawTransactionResult, error) {

	if params.Slippage == 0 {
		params.Slippage = 0.1
	}

	req := &jupiter.QuoteRequest{
		InputMint:   params.InputMint(),
		OutputMint:  params.OutputMint(),
		Amount:      conv.Str(params.Amount),
		SlippageBps: int(params.Slippage * 10000),
	}

	//if params.Fee != nil {
	//	req.PlatformFeeBps = params.Fee.Bps
	//}

	quote, err := t.jc.Quote(ctx, req)
	if err != nil {
		return nil, err
	}

	var prioritizationFeeLamports jupiter.PrioritizationFeeLamports

	if params.Priority == nil {
		if err = prioritizationFeeLamports.UnmarshalJSON([]byte(`"auto"`)); err != nil {
			return nil, err
		}

	} else {
		if err = prioritizationFeeLamports.UnmarshalJSON(conv.M2B(map[string]any{
			"priorityLevelWithMaxLamports": map[string]any{
				"priorityLevel": params.Priority.Level,
				"maxLamports":   params.Priority.Fee * math.Pow10(9),
			},
		})); err != nil {
			return nil, err
		}
	}

	var computeUnitPriceMicroLamports jupiter.ComputeUnitPriceMicroLamports
	if err = computeUnitPriceMicroLamports.UnmarshalJSON([]byte(`"auto"`)); err != nil {
		return nil, err
	}

	swapReq := &jupiter.SwapRequest{
		UserPublicKey:             params.Wallet,
		QuoteResponse:             *quote,
		DynamicComputeUnitLimit:   jupiter.Bool(true), // 加上这个之后 上链效率明显增加
		PrioritizationFeeLamports: &prioritizationFeeLamports,
		//ComputeUnitPriceMicroLamports: &computeUnitPriceMicroLamports,
		SkipUserAccountsRpcCalls: jupiter.Bool(true),
	}

	if params.Fee != nil {
		swapReq.FeeAccount = jupiter.String(params.Fee.Receiver)
	}

	swap, err := t.jc.Swap(ctx, swapReq)
	if err != nil {

		// todo 放到 jupiter内部
		if strings.Contains(err.Error(), "Missing token program") {
			return nil, jupiter.MissingTokenProgramError
		}

		return nil, err
	}

	tx := solana.MustTransactionFromBase64(swap.SwapTransaction)

	return &CreateRawTransactionResult{
		RawTX: tx,
	}, nil
}

func (t JupiterTrader) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error) {

	//if params.Slippage == 0 {
	//	params.Slippage = 0.1
	//}
	//
	//req := &jupiter.QuoteRequest{
	//	InputMint:   params.InputMint(),
	//	OutputMint:  params.OutputMint(),
	//	Amount:      conv.Str(params.Amount),
	//	SlippageBps: int(params.Slippage * 10000),
	//}
	//
	////if params.Fee != nil {
	////	req.PlatformFeeBps = params.Fee.Bps
	////}
	//
	//quote, err := t.jc.Quote(ctx, req)
	//if err != nil {
	//	return nil, err
	//}
	//
	//var prioritizationFeeLamports jupiter.PrioritizationFeeLamports
	//
	//if params.Priority == nil {
	//	if err = prioritizationFeeLamports.UnmarshalJSON([]byte(`"auto"`)); err != nil {
	//		return nil, err
	//	}
	//
	//} else {
	//	if err = prioritizationFeeLamports.UnmarshalJSON(conv.M2B(map[string]any{
	//		"priorityLevelWithMaxLamports": map[string]any{
	//			"priorityLevel": params.Priority.Level,
	//			"maxLamports":   params.Priority.Fee * math.Pow10(9),
	//		},
	//	})); err != nil {
	//		return nil, err
	//	}
	//}
	//
	//var computeUnitPriceMicroLamports jupiter.ComputeUnitPriceMicroLamports
	//if err = computeUnitPriceMicroLamports.UnmarshalJSON([]byte(`"auto"`)); err != nil {
	//	return nil, err
	//}
	//
	//swapReq := &jupiter.SwapRequest{
	//	UserPublicKey:             params.Wallet,
	//	QuoteResponse:             *quote,
	//	DynamicComputeUnitLimit:   jupiter.Bool(true), // 加上这个之后 上链效率明显增加
	//	PrioritizationFeeLamports: &prioritizationFeeLamports,
	//	//ComputeUnitPriceMicroLamports: &computeUnitPriceMicroLamports,
	//	SkipUserAccountsRpcCalls: jupiter.Bool(true),
	//}
	//
	//if params.Fee != nil {
	//	swapReq.FeeAccount = jupiter.String(params.Fee.Receiver)
	//}
	//
	//swap, err := t.jc.Swap(ctx, swapReq)
	//if err != nil {
	//
	//	// todo 放到 jupiter内部
	//	if strings.Contains(err.Error(), "Missing token program") {
	//		return nil, jupiter.MissingTokenProgramError
	//	}
	//
	//	return nil, err
	//}
	//
	//tx := solana.MustTransactionFromBase64(swap.SwapTransaction)
	//
	//return &CreateTransactionResult{
	//	TxHash:    sign,
	//	Confirmed: ok,
	//	Quote: &Quote{
	//		Token:          params.Token,
	//		TokenAmountRaw: helper.Select(params.IsBuy(), quote.OutAmount, quote.InAmount),
	//		SolAmountRaw:   helper.Select(params.IsBuy(), quote.InAmount, quote.OutAmount),
	//		Side:           params.Side,
	//	},
	//}, nil

	return nil, nil
}

func NewJupiterTrader(solana *rpc.Client, jc *jupiter.SwapClient) ITrader {
	return &JupiterTrader{solana: solana, jc: jc}

}
