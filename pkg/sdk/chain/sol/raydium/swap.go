package raydium

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/conv"
)

type SwapParams struct {
	ComputeUnitPriceMicroLamports int64 `json:"computeUnitPriceMicroLamports"`
	SwapResponse                  *QuoteResult
	TxVersion                     string // V0 or LEGACY
	Wallet                        string
}

type SwapResult struct {
	Id      string `json:"id"`
	Version string `json:"version"`
	Success bool   `json:"success"`
	Data    []Swap `json:"data"`
}

type Swap struct {
	Transaction string `json:"transaction"`
}

func (t *Client) Swap(ctx context.Context, params SwapParams) (*Swap, error) {
	//
	//address, u, err := solana.FindAssociatedTokenAddress(
	//	solana.MPK("3Cyrb9NABtEMJWv5yJ5pnccJPt3UYb1TmyiACu98NQts"),
	//	solana.MPK("Ho6wN4ff7RdTdXE1UsCZjrjuFVMHyRFTv1oBdbSECnJS"),
	//)
	//
	//fmt.Println(address.String(), u, err)

	inputSol, outputSol := params.SwapResponse.Data.InputMint == solana.SolMintString, params.SwapResponse.Data.OutputMint == solana.SolMintString

	body := map[string]interface{}{
		"computeUnitPriceMicroLamports": conv.Str(params.ComputeUnitPriceMicroLamports),
		"swapResponse":                  params.SwapResponse,
		"txVersion":                     params.TxVersion,
		"wallet":                        params.Wallet,
		"wrapSol":                       inputSol,
		"unwrapSol":                     outputSol,
	}

	if inputSol {
		ata, _, _ := solana.FindAssociatedTokenAddress(
			solana.MPK(params.Wallet),
			solana.MPK(params.SwapResponse.Data.OutputMint),
		)
		body["outputAccount"] = ata.String()
	}

	if outputSol {
		ata, _, _ := solana.FindAssociatedTokenAddress(
			solana.MPK(params.Wallet),
			solana.MPK(params.SwapResponse.Data.InputMint),
		)
		body["inputAccount"] = ata.String()
	}

	result, err := resty.New().R().SetContext(ctx).
		SetBody(body).
		Post(t.endpoint + "/transaction/swap-base-in")

	if err != nil {
		return nil, err
	}

	if result.StatusCode() != 200 {
		return nil, errors.New(result.String())
	}

	var rsp SwapResult

	err = json.Unmarshal(result.Body(), &rsp)
	if err != nil {
		return nil, err
	}

	if len(rsp.Data) == 0 {
		return nil, errors.New("swap result empty")
	}

	return &rsp.Data[0], nil
}
