package trader

import (
	"context"
	"math"
	"math/big"
	"store/pkg/sdk/chain/sol/pumpdotfunsdk"
	"store/pkg/sdk/chain/sol/quicknode"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/chain/sol/solana/programs/system"
	"store/pkg/sdk/chain/sol/solana/rpc"
	"store/pkg/sdk/chain/sol/solana/rpc/ws"
	"store/pkg/sdk/helper/mathz"
)

func NewPumpfunInstructionsTrader(solana *rpc.Client, solanaWS *ws.Client, qn *quicknode.Client,
) ITrader {
	return &PumpfunInstructionsTrader{
		solana:   solana,
		solanaWS: solanaWS,
		qn:       qn,
	}
}

type PumpfunInstructionsTrader struct {
	solana   *rpc.Client
	solanaWS *ws.Client
	qn       *quicknode.Client
}

func (t *PumpfunInstructionsTrader) CreateRawTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateRawTransactionResult, error) {

	slippage := 0.3
	if params.Slippage > 0 {
		slippage = params.Slippage
	}

	wallet := solana.MPK(params.Wallet)
	token := solana.MPK(params.Token)

	var instructions []solana.Instruction
	var toAmount *big.Int
	var err error

	//fees, err := t.qn.EstimatePriorityFees(ctx, &quicknode.EstimatePriorityFeesParams{
	//	LastNBlocks: 0,
	//	Account:     solana.ProgramPumpFun.String(),
	//})
	//if err != nil {
	//	return nil, err
	//}

	if params.Side == "buy" {
		instructions, toAmount, err = pumpdotfunsdk.CreateBuyInstructions(
			t.solana,
			wallet,
			token,
			uint64(params.Amount),
			uint(mathz.BasicPoints(slippage)),
			uint64(params.Priority.Fee*math.Pow10(9)),
			//uint64(fees.PerComputeUnit.High),
		)
	} else {
		instructions, toAmount, err = pumpdotfunsdk.CreateSellInstructions(
			t.solana,
			wallet,
			token,
			uint64(params.Amount),
			uint(mathz.BasicPoints(slippage)),
			uint64(params.Priority.Fee*math.Pow10(9)),
			//uint64(fees.PerComputeUnit.High),
		)
	}

	if err != nil {
		return nil, err
	}

	// 构建交易
	builder := solana.NewTransactionBuilder()

	for _, x := range instructions {
		builder.AddInstruction(x)
	}

	// 确认手续费
	if params.Fee != nil {

		var feeAmount float64

		if params.IsBuy() {
			feeAmount = float64(params.Amount) * params.Fee.Rate
		} else {
			feeAmount = float64(toAmount.Int64()) * params.Fee.Rate
		}

		builder.AddInstruction(system.NewTransferInstruction(
			uint64(feeAmount),
			wallet,
			solana.MustPublicKeyFromBase58(params.Fee.Receiver),
		).Build())
	}

	builder.SetFeePayer(solana.MPK(params.Wallet))

	tx, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return &CreateRawTransactionResult{
		RawTX: tx,
	}, nil
}

func (t *PumpfunInstructionsTrader) CreateInstructions(ctx context.Context, params *CreateCreateInstructionsParams) (*CreateInstructionsResult, error) {
	if params.Slippage == 0 {
		params.Slippage = 0.05
	}

	wallet := solana.MPK(params.Wallet)
	token := solana.MPK(params.Token)

	var instructions []solana.Instruction
	var toAmount *big.Int
	var err error

	var amount, quoteAmount int64

	if params.Side == "buy" {
		instructions, toAmount, err = pumpdotfunsdk.CreateBuyInstructions(
			t.solana,
			wallet,
			token,
			uint64(params.Amount),
			uint(mathz.BasicPoints(params.Slippage)),
			uint64(params.Priority.Fee*math.Pow10(9)),
		)

		if err != nil {
			return nil, err
		}

		quoteAmount = params.Amount
		amount = toAmount.Int64()

	} else {
		instructions, toAmount, err = pumpdotfunsdk.CreateSellInstructions(
			t.solana,
			wallet,
			token,
			uint64(params.Amount),
			uint(mathz.BasicPoints(params.Slippage)),
			uint64(params.Priority.Fee*math.Pow10(9)),
		)

		if err != nil {
			return nil, err
		}

		amount = params.Amount
		quoteAmount = toAmount.Int64()
	}

	return &CreateInstructionsResult{
		Instructions: instructions,
		Quotation: &Quotation{
			Token:       params.Token,
			Amount:      amount,
			QuoteAmount: quoteAmount,
			Side:        params.Side,
		},
	}, nil
}

func (t *PumpfunInstructionsTrader) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error) {
	//TODO implement me
	panic("implement me")
}
