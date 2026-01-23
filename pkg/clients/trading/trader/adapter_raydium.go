package trader

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"store/pkg/sdk/chain/sol/raydium"
	"store/pkg/sdk/chain/sol/solana"
	addresslookuptable "store/pkg/sdk/chain/sol/solana/programs/address-lookup-table"
	"store/pkg/sdk/chain/sol/solana/rpc"
)

type RaydiumTrader struct {
	solana *rpc.Client
}

func (r RaydiumTrader) CreateRawTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateRawTransactionResult, error) {
	c := raydium.NewClient()

	var slippage = 0.05
	if params.Slippage > 0 {
		slippage = params.Slippage
	}

	quote, err := c.Quote(ctx, raydium.QuoteParams{
		InputMint:  params.InputMint(),
		OutputMint: params.OutputMint(),
		Amount:     params.Amount,
		Slippage:   slippage,
		TxVersion:  "V0",
	})
	if err != nil {
		return nil, err
	}

	swap, err := c.Swap(ctx, raydium.SwapParams{
		//ComputeUnitPriceMicroLamports: 10000000,
		SwapResponse: quote,
		TxVersion:    "V0",
		Wallet:       params.Wallet,
	})
	if err != nil {
		return nil, err
	}

	tx := solana.MustTransactionFromBase64(swap.Transaction)

	return &CreateRawTransactionResult{
		RawTX: tx,
	}, nil
}

func (r RaydiumTrader) CreateInstructions(ctx context.Context, params *CreateCreateInstructionsParams) (*CreateInstructionsResult, error) {
	c := raydium.NewClient()

	var slippage = 0.05
	if params.Slippage > 0 {
		slippage = params.Slippage
	}

	quote, err := c.Quote(ctx, raydium.QuoteParams{
		InputMint:  params.InputMint(),
		OutputMint: params.OutputMint(),
		Amount:     params.Amount,
		Slippage:   slippage,
	})
	if err != nil {
		return nil, err
	}

	swap, err := c.Swap(ctx, raydium.SwapParams{
		ComputeUnitPriceMicroLamports: 10000000,
		SwapResponse:                  quote,
		TxVersion:                     "V0",
		Wallet:                        params.Wallet,
	})
	if err != nil {
		return nil, err
	}

	tx := solana.MustTransactionFromBase64(swap.Transaction)

	spew.Dump(tx)

	if len(tx.Message.AddressTableLookups) > 0 {
		table, err := addresslookuptable.GetAddressLookupTable(ctx, r.solana, tx.Message.AddressTableLookups[0].AccountKey)
		if err != nil {
			return nil, err
		}

		fmt.Println(table)
	}

	var instructions []solana.Instruction
	for _, x := range tx.Message.Instructions {

		var accounts solana.AccountMetaSlice
		for _, xx := range x.Accounts {

			account, _ := tx.Message.Account(xx)

			accounts = append(accounts, &solana.AccountMeta{
				PublicKey:  account,
				IsWritable: true,
				IsSigner:   account.String() == params.Wallet,
			})
		}

		//data, _ := base58.Decode(string(x.Data))

		programAccount, _ := tx.Message.Account(x.ProgramIDIndex)

		y := solana.NewInstruction(
			programAccount,
			accounts,
			x.Data,
		)

		instructions = append(instructions, y)

	}

	return &CreateInstructionsResult{
		Instructions: instructions,
	}, nil
}

func (r RaydiumTrader) CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error) {
	//TODO implement me
	panic("implement me")
}

func NewRaydiumTrader(solana *rpc.Client) ITrader {
	return &RaydiumTrader{solana: solana}
}
