package pumpdotfunsdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"store/pkg/sdk/chain/sol/pumpdotfunsdk/pump"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/chain/sol/solana/rpc"

	// General solana packages.
	confirm "store/pkg/sdk/chain/sol/solana/rpc/sendAndConfirmTransaction"
	"store/pkg/sdk/chain/sol/solana/rpc/ws"

	// This package interacts with the Compute Budget program, allowing
	// to easily get instruction to set compute budget limit/price for example.
	cb "store/pkg/sdk/chain/sol/solana/programs/compute-budget"
	// This package interacts with the Solana system program, allowing
	// to transfer solana for example.
	"store/pkg/sdk/chain/sol/solana/programs/system"
	// This package interacts with the Token program, allowing
	// to create a token for example.
	"store/pkg/sdk/chain/sol/solana/programs/token"
	// This package interacts with the Associated Token Account program
	// allowing to create/close an associated token account for example.
	associatedtokenaccount "store/pkg/sdk/chain/sol/solana/programs/associated-token-account"
)

// Contains commonly used addresses with the pump.fun program, that are not present
// in the generated code, from its IDL file.
var (
	// Global account address for pump.fun
	globalPumpFunAddress = solana.MustPublicKeyFromBase58("4wTV1YmiEkRvAtNtsSGPtUrqRYQMe5SKy2uB4Jjaxnjf")
	// Pump.fun mint authority
	pumpFunMintAuthority = solana.MustPublicKeyFromBase58("TSLvdd1pWpHVjahSpsvCXUbgwsL3JAcvokwaKt1eokM")
	// Pump.fun event authority
	pumpFunEventAuthority = solana.MustPublicKeyFromBase58("Ce6TQqeHC9p8KetsN6JsjHK7UTZk7nasjjnr7XxXp9F1")
	// Pump.fun fee recipient
	pumpFunFeeRecipient = solana.MustPublicKeyFromBase58("CebN5WGQ4jvEPvsVU4EoHEpgzq1VV7AbicfhtW4xC9iM")
)

// SetDevnetMode sets the pump.fun program addresses to the devnet addresses.
// It is important to call this function if you are using the devnet.
func SetDevnetMode() {
	// This is the address you want to use as pump.fun fee recipient on devnet, otherwise, it
	// will not work, as the official pump.fun fee recipient account is not initialized on devnet.
	// I know, using global variables is ugly, but passing this address around everywhere
	// (in BuyToken / SellToken), while it's actually a constant on mainnet is even uglier,
	// considering that there is no other difference.
	pumpFunFeeRecipient = solana.MustPublicKeyFromBase58("68yFSZxzLWJXkxxRGydZ63C6mHx1NLEDWmwN9Lb5yySg")
}

type BondingCurvePublicKeys struct {
	BondingCurve           solana.PublicKey
	AssociatedBondingCurve solana.PublicKey
}

// getBondingCurveAndAssociatedBondingCurve returns the bonding curve and associated bonding curve, in a structured format.
func getBondingCurveAndAssociatedBondingCurve(mint solana.PublicKey) (*BondingCurvePublicKeys, error) {
	// Derive bonding curve address.
	// define the seeds used to derive the PDA
	// getProgramDerivedAddress equivalent.
	seeds := [][]byte{
		[]byte("bonding-curve"),
		mint.Bytes(),
	}
	bondingCurve, _, err := solana.FindProgramAddress(seeds, pump.ProgramID)
	if err != nil {
		return nil, fmt.Errorf("failed to derive bonding curve address: %w", err)
	}
	// Derive associated bonding curve address.
	associatedBondingCurve, _, err := solana.FindAssociatedTokenAddress(
		bondingCurve,
		mint,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to derive associated bonding curve address: %w", err)
	}
	return &BondingCurvePublicKeys{
		BondingCurve:           bondingCurve,
		AssociatedBondingCurve: associatedBondingCurve,
	}, nil
}

func getComputUnitPriceInstr(rpcClient *rpc.Client, user solana.PrivateKey) (*cb.SetComputeUnitPrice, error) {
	// create priority fee instructions
	out, err := rpcClient.GetRecentPrioritizationFees(context.TODO(), solana.PublicKeySlice{user.PublicKey(), pump.ProgramID, pumpFunMintAuthority, globalPumpFunAddress, solana.TokenMetadataProgramID, system.ProgramID, token.ProgramID, associatedtokenaccount.ProgramID, solana.SysVarRentPubkey, pumpFunEventAuthority})
	if err != nil {
		return nil, fmt.Errorf("failed to get recent prioritization fees: %w", err)
	}
	var median uint64
	length := uint64(len(out))
	for _, fee := range out {
		median = fee.PrioritizationFee
	}
	median /= length
	cupInst := cb.NewSetComputeUnitPriceInstruction(median)
	return cupInst, nil
}

func CreateToken(rpcClient *rpc.Client, wsClient *ws.Client, user solana.PrivateKey, mint *solana.Wallet, name string, symbol string, uri string, buyAmountLamports uint64, slippageBasisPoint uint) (string, error) {
	bondingCurveData, err := getBondingCurveAndAssociatedBondingCurve(mint.PublicKey())
	if err != nil {
		return "", fmt.Errorf("failed to get bonding curve and associated bonding curve: %w", err)
	}
	// Get token metadata address
	metadata, _, err := solana.FindTokenMetadataAddress(mint.PublicKey())
	if err != nil {
		return "", fmt.Errorf("can't find token metadata address: %w", err)
	}

	// Default pump.fun compute limit is 250k, so we set the same here.
	culInst := cb.NewSetComputeUnitLimitInstruction(uint32(250000))
	cupInst, err := getComputUnitPriceInstr(rpcClient, user)
	if err != nil {
		return "", fmt.Errorf("failed to get compute unit price instructions: %w", err)
	}
	// Create the pump fun instruction
	instr := pump.NewCreateInstruction(
		name,
		symbol,
		uri,
		mint.PublicKey(),
		pumpFunMintAuthority,
		bondingCurveData.BondingCurve,
		bondingCurveData.AssociatedBondingCurve,
		globalPumpFunAddress,
		solana.TokenMetadataProgramID,
		metadata,
		user.PublicKey(),
		system.ProgramID,
		token.ProgramID,
		associatedtokenaccount.ProgramID,
		solana.SysVarRentPubkey,
		pumpFunEventAuthority,
		pump.ProgramID,
	)
	instruction := instr.Build()
	// get recent block hash
	recent, err := rpcClient.GetLatestBlockhash(context.TODO(), rpc.CommitmentFinalized)
	if err != nil {
		return "", fmt.Errorf("error while getting recent block hash: %w", err)
	}
	instructions := []solana.Instruction{
		culInst.Build(),
		cupInst.Build(),
		instruction,
	}
	// get buy instructions
	if buyAmountLamports > 0 {
		buyInstructions, _, err := getBuyInstructions(rpcClient, mint.PublicKey(), user.PublicKey(), buyAmountLamports, slippageBasisPoint)
		if err != nil {
			return "", fmt.Errorf("failed to get buy instructions: %w", err)
		}
		instructions = append(instructions, buyInstructions...)
	}
	tx, err := solana.NewTransaction(
		instructions,
		recent.Value.Blockhash,
		solana.TransactionPayer(user.PublicKey()),
	)
	if err != nil {
		return "", fmt.Errorf("error while creating new transaction: %w", err)
	}
	_, err = tx.Sign(
		func(key solana.PublicKey) *solana.PrivateKey {
			if user.PublicKey().Equals(key) {
				return &user
			}
			if mint.PublicKey().Equals(key) {
				return &mint.PrivateKey
			}
			return nil
		},
	)
	if err != nil {
		return "", fmt.Errorf("can't sign transaction: %w", err)
	}
	// Send transaction, and wait for confirmation:
	sig, _, err := confirm.SendAndConfirmTransaction(
		context.TODO(),
		rpcClient,
		wsClient,
		tx,
	)
	if err != nil {
		return "", fmt.Errorf("can't send and confirm new transaction: %w", err)
	}
	return sig.String(), nil
}

type CreateTokenMetadataRequest struct {
	Filename    string
	Name        string
	Symbol      string
	Description string
	Twitter     string
	Telegram    string
	Website     string
}

type CreateTokenMetadataResponse struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	ShowName    bool   `json:"showName"`
	CreatedOn   string `json:"createdOn"`
	Twitter     string `json:"twitter"`
	Telegram    string `json:"telegram"`
	Website     string `json:"website"`

	Image       string `json:"image"`
	MetadataUri string `json:"metadataUri"`
}

func CreateTokenMetadata(client *http.Client, create CreateTokenMetadataRequest) (*CreateTokenMetadataResponse, error) {
	// Create a buffer to hold the form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Add the file from URL
	resp, err := http.Get(create.Filename)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Create the form file
	part, err := writer.CreateFormFile("file", "image.png")
	if err != nil {
		return nil, err
	}
	// Copy the file content to the form file
	_, err = io.Copy(part, resp.Body)
	if err != nil {
		return nil, err
	}

	// Add the other form fields
	writer.WriteField("name", create.Name)
	writer.WriteField("symbol", create.Symbol)
	writer.WriteField("description", create.Description)
	writer.WriteField("twitter", create.Twitter)
	writer.WriteField("telegram", create.Telegram)
	writer.WriteField("website", create.Website)
	writer.WriteField("showName", "true")

	// Close the writer to finalize the form data
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", "https://pump.fun/api/ipfs", &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the HTTP request
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the JSON response
	var result CreateTokenMetadataResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
