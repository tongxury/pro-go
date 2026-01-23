package trader

import (
	"context"
	"encoding/base64"
	"store/pkg/sdk/chain/sol/jupiter"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/chain/sol/solana/rpc"
)

type JupiterInstructionsTrader struct {
	solana *rpc.Client
	jc     *jupiter.SwapClient
}

func (t JupiterInstructionsTrader) CreateInstructions(ctx context.Context, params *CreateCreateInstructionsParams) (*CreateInstructionsResult, error) {

	//TODO implement me
	panic("implement me")
}

func (t JupiterInstructionsTrader) CreateRawTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateRawTransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func (t JupiterInstructionsTrader) getFeeInstruction() (*solana.Instruction, error) {

	////
	////if !feeAmount.GtZero() {
	////	builder.AddInstruction(system.NewTransferInstruction(
	////		feeAmount.,
	////		payer,
	////		solana.MustPublicKeyFromBase58("3xUit8UR1cEtKGtUU9dW8ViHioRi94CsfCpxiwJdWCWz"),
	////	).Build())
	////
	////}
	//
	//var sourcePublicKey solana.PublicKey
	//for _, x := range swapInstructions.SetupInstructions {
	//	if x.ProgramId == solana.SPLAssociatedTokenAccountProgramID.String() {
	//
	//		for _, xx := range x.Accounts {
	//
	//			if !helper.Contains(
	//				[]string{
	//					params.Wallet,
	//					solana.SPLAssociatedTokenAccountProgramID.String(),
	//					solana.SolMint.String(),
	//					solana.SystemProgramID.String(),
	//					solana.TokenProgramID.String(),
	//				},
	//				xx.Pubkey,
	//			) {
	//				sourcePublicKey = solana.MustPublicKeyFromBase58(xx.Pubkey)
	//			}
	//		}
	//	}
	//}
	//
	//sourcePublicKey.IsZero()
	//
	//builder.AddInstruction(
	//	token.NewTransferInstruction(
	//		1472222,
	//		solana.MustPublicKeyFromBase58("7qaqheEJuLDWNnPhdMBfF3i3HSbQCGmYmsp3rJX6z7dz"),
	//		solana.MustPublicKeyFromBase58("2YWbZ9MgzX6R3YTTqapYGEA5rNt5gbb6FQUtQ5fsqacd"),
	//		//solana.MustPublicKeyFromBase58("7qaqheEJuLDWNnPhdMBfF3i3HSbQCGmYmsp3rJX6z7dz"),
	//		//solana.MustPublicKeyFromBase58(params.Fee.Receiver),
	//		payer,
	//
	//		[]solana.PublicKey{},
	//	).Build())

	return nil, nil
}

func (t JupiterInstructionsTrader) asSolanaInstruction(instruction jupiter.Instruction) solana.Instruction {

	var accounts solana.AccountMetaSlice
	for _, x := range instruction.Accounts {
		accounts = append(accounts, &solana.AccountMeta{
			PublicKey:  solana.MustPublicKeyFromBase58(x.Pubkey),
			IsWritable: x.IsWritable,
			IsSigner:   x.IsSigner,
		})
	}

	data, _ := base64.StdEncoding.DecodeString(instruction.Data)

	return solana.NewInstruction(
		solana.MustPublicKeyFromBase58(instruction.ProgramId),
		accounts,
		data,
	)

}

func (t JupiterInstructionsTrader) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error) {

	//if params.Slippage == 0 {
	//	params.Slippage = 0.1
	//}
	//
	//var buyAmount, feeAmount = params.Amount, bigdecimal.Zero
	//// 确认手续费
	//if params.Fee != nil {
	//	feeAmount = params.Amount.MulFloat64(params.Fee.Rate)
	//
	//	if params.IsBuy() {
	//		// 买币预先扣除手续费
	//		buyAmount = params.Amount.Sub(feeAmount)
	//	} else {
	//		buyAmount = params.Amount
	//	}
	//}
	//
	//req := &jupiter.QuoteRequest{
	//	InputMint:   params.InputMint(),
	//	OutputMint:  params.OutputMint(),
	//	Amount:      buyAmount.Raw(),
	//	SlippageBps: int(params.Slippage * 10000),
	//}
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
	//swapInstructions, err := t.jc.SwapInstructions(ctx, swapReq)
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
	//// 构建tx
	//builder := solana.NewTransactionBuilder()
	//
	//payer := solana.MustPublicKeyFromBase58(params.Wallet)
	////feeReceiver := solana.MustPublicKeyFromBase58("3xUit8UR1cEtKGtUU9dW8ViHioRi94CsfCpxiwJdWCWz")
	//
	//if swapInstructions.TokenLedgerInstruction != nil {
	//	builder = builder.AddInstruction(t.asSolanaInstruction(*swapInstructions.TokenLedgerInstruction))
	//}
	//
	//for _, x := range swapInstructions.ComputeBudgetInstructions {
	//	builder = builder.AddInstruction(t.asSolanaInstruction(x))
	//}
	//
	//for _, x := range swapInstructions.SetupInstructions {
	//	builder = builder.AddInstruction(t.asSolanaInstruction(x))
	//}
	//
	//if swapInstructions.SwapInstruction != nil {
	//	builder = builder.AddInstruction(t.asSolanaInstruction(*swapInstructions.SwapInstruction))
	//}
	//
	//if swapInstructions.CleanupInstruction != nil {
	//	builder = builder.AddInstruction(t.asSolanaInstruction(*swapInstructions.CleanupInstruction))
	//}
	//
	//for _, x := range swapInstructions.OtherInstructions {
	//	builder = builder.AddInstruction(t.asSolanaInstruction(x))
	//}
	//
	//latestBlockHashResult, err := t.solana.GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	//if err != nil {
	//	return nil, err
	//}
	//
	//builder = builder.SetRecentBlockHash(latestBlockHashResult.Value.Blockhash)
	//builder = builder.SetFeePayer(payer)
	//
	//tx, err := builder.Build()
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//sign, ok, err := t.solana.SendTransactionAndWaitUtilConfirmed(ctx,
	//	rpc.SendTransactionParams{
	//		LatestBlock: &rpc.LatestBlock{
	//			Hash: latestBlockHashResult.Value.Blockhash,
	//			Slot: latestBlockHashResult.Value.LastValidBlockHeight,
	//		},
	//		Transaction: tx,
	//		PrivateKeys: []string{params.PrivateKey},
	//	})
	//
	//if err != nil {
	//
	//	if strings.Contains(err.Error(), "SlippageToleranceExceeded") {
	//		return nil, jupiter.SlippageToleranceExceededError
	//	}
	//
	//	return nil, fmt.Errorf("send transaction err: %v", err)
	//}

	return &CreateTransactionResult{
		//TxHash:    sign,
		//Confirmed: ok,
		//Quote: &Quote{
		//	Token:          params.Token,
		//	TokenAmountRaw: helper.Select(params.IsBuy(), quote.OutAmount, quote.InAmount),
		//	SolAmountRaw:   helper.Select(params.IsBuy(), quote.InAmount, quote.OutAmount),
		//	Side:           params.Side,
		//},
	}, nil
}

func NewJupiterInstructionsTrader(solana *rpc.Client, jc *jupiter.SwapClient) ITrader {
	return &JupiterInstructionsTrader{solana: solana, jc: jc}

}
