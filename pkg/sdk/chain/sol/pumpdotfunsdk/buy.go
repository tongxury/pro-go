package pumpdotfunsdk

import (
	"context"
	"fmt"
	"math/big"
	"store/pkg/sdk/chain/sol/pumpdotfunsdk/pump"
	"store/pkg/sdk/chain/sol/solana"
	associatedtokenaccount "store/pkg/sdk/chain/sol/solana/programs/associated-token-account"
	cb "store/pkg/sdk/chain/sol/solana/programs/compute-budget"
	"store/pkg/sdk/chain/sol/solana/programs/system"
	"store/pkg/sdk/chain/sol/solana/programs/token"
	"store/pkg/sdk/chain/sol/solana/rpc"
)

// checks if the associated token account for the mint and our bot's public key exists.
func shouldCreateAta(rpcClient *rpc.Client, ata solana.PublicKey) (bool, error) {
	_, err := rpcClient.GetAccountInfo(context.TODO(), ata)
	if err == nil {
		return false, nil
	}
	return true, nil
}

func buildComputeUnitInstructions(priorityFee uint64) []solana.Instruction {

	var unitLimit = 2500000
	// create priority fee instructions
	culInst := cb.NewSetComputeUnitLimitInstruction(uint32(unitLimit))

	microLamports := uint64(100000)
	if priorityFee > 0 {
		microLamports = priorityFee * (1000000) / uint64(unitLimit)
	}

	if microLamports > 10000000 {
		microLamports = 10000000
	}

	cupInst := cb.NewSetComputeUnitPriceInstruction(microLamports)

	//culInst := cb.NewSetComputeUnitLimitInstruction(uint32(250000))
	//cupInst := cb.NewSetComputeUnitPriceInstruction(100000)

	return []solana.Instruction{
		cupInst.Build(),
		culInst.Build(),
	}
}

// buyToken buys a token from the bonding curve.
// The amount is the amount of tokens to buy, and the sol is the amount of SOL to pay.
// The mintAddr is the address of the mint of the token.
// This function will send a transaction to the network to buy the token.
// This function will return an error if the transaction fails.
func CreateBuyInstructions(
	rpcClient *rpc.Client,
	wallet solana.PublicKey,
	mint solana.PublicKey,
	buyAmountLamports uint64,
	slippageBasisPoint uint,
	priorityFee uint64,
) ([]solana.Instruction, *big.Int, error) {

	instructions := buildComputeUnitInstructions(priorityFee)
	// get buy instructions
	buyInstructions, amount, err := getBuyInstructions(
		rpcClient,
		mint,
		wallet,
		buyAmountLamports,
		slippageBasisPoint,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get buy instructions: %w", err)
	}
	instructions = append(instructions, buyInstructions...)

	return instructions, amount, err
	//// get recent block hash
	//recent, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	//if err != nil {
	//	return "", fmt.Errorf("error while getting recent block hash: %w", err)
	//}
	//// create new transaction
	//tx, err := solana.NewTransaction(
	//	instructions,
	//	recent.Value.Blockhash,
	//	solana.TransactionPayer(wallet.PublicKey()),
	//)
	//if err != nil {
	//	return "", fmt.Errorf("error while creating new transaction: %w", err)
	//}
	//_, err = tx.Sign(
	//	func(key solana.PublicKey) *solana.PrivateKey {
	//		if wallet.PublicKey().Equals(key) {
	//			return &wallet
	//		}
	//		return nil
	//	},
	//)
	//if err != nil {
	//	return "", fmt.Errorf("can't sign transaction: %w", err)
	//}
	//
	//// Send transaction, and wait for confirmation:
	//
	//sig, err := rpcClient.SendTransaction(context.Background(), tx)
	//
	////sig, err := confirm.SendAndConfirmTransaction(
	////	context.TODO(),
	////	rpcClient,
	////	wsClient,
	////	tx,
	////)
	//if err != nil {
	//	return "", fmt.Errorf("can't send and confirm new transaction: %w", err)
	//}
	//return sig.String(), nil
}

func getBuyInstructions(
	rpcClient *rpc.Client,
	mint solana.PublicKey,
	user solana.PublicKey,
	solAmount uint64,
	slippageBasisPoint uint,
) ([]solana.Instruction, *big.Int, error) {
	bondingCurveData, err := getBondingCurveAndAssociatedBondingCurve(mint)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get bonding curve data: %w", err)
	}
	// NOTE: buy transaction for the token
	var instructions []solana.Instruction
	ata, _, err := solana.FindAssociatedTokenAddress(
		user,
		mint,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to derive associated token account: %w", err)
	}
	shouldCreateATA, err := shouldCreateAta(rpcClient, ata)
	if err != nil {
		return nil, nil, fmt.Errorf("can't check if we should create ATA: %w", err)
	}
	if shouldCreateATA {
		ataInstr, err := associatedtokenaccount.NewCreateInstruction(user, user, mint).
			ValidateAndBuild()
		if err != nil {
			return nil, nil, fmt.Errorf("can't create associated token account: %w", err)
		}
		instructions = append(instructions, ataInstr)
	}

	bondingCurve, err := fetchBondingCurve(rpcClient, bondingCurveData.BondingCurve)
	if err != nil {
		return nil, nil, fmt.Errorf("can't fetch bonding curve: %w", err)
	}
	// We set 2% slippage.
	percentage := convertSlippageBasisPointsToPercentage(slippageBasisPoint)
	buy := calculateBuyQuote(solAmount, bondingCurve, percentage)
	buyInstr := pump.NewBuyInstruction(
		buy.Uint64(),
		solAmount,
		globalPumpFunAddress,
		pumpFunFeeRecipient,
		mint,
		bondingCurveData.BondingCurve,
		bondingCurveData.AssociatedBondingCurve,
		ata,
		user,
		system.ProgramID,
		token.ProgramID,
		solana.SysVarRentPubkey,
		pumpFunEventAuthority,
		pump.ProgramID,
	)
	buyInstruction := buyInstr.Build()
	instructions = append(instructions, buyInstruction)
	return instructions, buy, nil
}

func convertSlippageBasisPointsToPercentage(slippageBasisPoint uint) float64 {
	return 1.0 - float64(slippageBasisPoint)/10e3
}

// calculateBuyQuote calculates how many tokens can be purchased given a specific amount of SOL, bonding curve data, and percentage.
// solAmount is the amount of sol that you want to buy
// bondingCurve is the BondingCurveData, that includes the real, virtual token/sol reserves, in order to calculate the price.
// percentage is what you want to use to set the slippage. For 2% slippage, you want to set the percentage to 0.98.
func calculateBuyQuote(
	solAmount uint64,
	bondingCurve *BondingCurveData,
	percentage float64,
) *big.Int {
	// Convert solAmount to *big.Int
	solAmountBig := big.NewInt(int64(solAmount))

	// Clone bonding curve data to avoid mutations
	virtualSolReserves := new(big.Int).Set(bondingCurve.VirtualSolReserves)
	virtualTokenReserves := new(big.Int).Set(bondingCurve.VirtualTokenReserves)

	// Compute the new virtual reserves
	newVirtualSolReserves := new(big.Int).Add(virtualSolReserves, solAmountBig)
	invariant := new(big.Int).Mul(virtualSolReserves, virtualTokenReserves)
	newVirtualTokenReserves := new(big.Int).Div(invariant, newVirtualSolReserves)

	// Calculate the tokens to buy
	tokensToBuy := new(big.Int).Sub(virtualTokenReserves, newVirtualTokenReserves)

	// Apply the percentage reduction (e.g., 95% or 0.95)
	// Convert the percentage to a multiplier (0.95) and apply to tokensToBuy
	percentageMultiplier := big.NewFloat(percentage)
	tokensToBuyFloat := new(big.Float).SetInt(tokensToBuy)
	finalTokens := new(big.Float).Mul(tokensToBuyFloat, percentageMultiplier)

	// Convert the result back to *big.Int
	finalTokensBig, _ := finalTokens.Int(nil)

	return finalTokensBig
}
