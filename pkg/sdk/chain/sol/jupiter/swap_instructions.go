package jupiter

import (
	"context"
	"errors"
	"store/pkg/sdk/helper/restyd"
)

type Instruction struct {
	ProgramId string   `json:"programId"`
	Accounts  Accounts `json:"accounts"`
	Data      string   `json:"data"`
}

type Account struct {
	Pubkey     string `json:"pubkey"`
	IsSigner   bool   `json:"isSigner"`
	IsWritable bool   `json:"isWritable"`
}
type Accounts []Account

func (ts Accounts) Pubkeys() []string {
	ret := make([]string, len(ts))
	for i, t := range ts {
		ret[i] = t.Pubkey
	}
	return ret

}

type SwapInstructionsResponse struct {
	TokenLedgerInstruction      *Instruction  `json:"tokenLedgerInstruction"`
	ComputeBudgetInstructions   []Instruction `json:"computeBudgetInstructions"`
	SetupInstructions           []Instruction `json:"setupInstructions"`
	SwapInstruction             *Instruction  `json:"swapInstruction"`
	CleanupInstruction          *Instruction  `json:"cleanupInstruction"`
	OtherInstructions           []Instruction `json:"otherInstructions"`
	AddressLookupTableAddresses []interface{} `json:"addressLookupTableAddresses"`
	PrioritizationFeeLamports   int           `json:"prioritizationFeeLamports"`
	ComputeUnitLimit            int           `json:"computeUnitLimit"`
	PrioritizationType          struct {
		ComputeBudget struct {
			MicroLamports          int `json:"microLamports"`
			EstimatedMicroLamports int `json:"estimatedMicroLamports"`
		} `json:"computeBudget"`
	} `json:"prioritizationType"`
	DynamicSlippageReport interface{} `json:"dynamicSlippageReport"`
	SimulationError       interface{} `json:"simulationError"`
}

func (t *SwapClient) SwapInstructions(ctx context.Context, params *SwapRequest) (*SwapInstructionsResponse, error) {

	url := t.endpoint + "/swap-instructions"

	swap, err := restyd.Post[SwapInstructionsResponse](ctx,
		restyd.Params{
			Url:  url,
			Body: params,
		})
	if err != nil {
		return nil, err
	}

	//
	//var swap SwapInstructionsResponse
	//
	//rsp, err := resty.New().R().
	//	SetContext(ctx).
	//	SetHeader("Content-Type", "application/json").
	//	SetBody(params).
	//	Post(url)
	//
	//if err != nil {
	//	return nil, err
	//}
	//
	//err = json.Unmarshal(rsp.Body(), &swap)
	//if err != nil {
	//	return nil, err
	//}

	//if rsp.StatusCode() != 200 {
	//	return nil, fmt.Errorf("SwapInstructions err, status: %d, body: %v", rsp.StatusCode(), rsp.String())
	//}

	if swap.SwapInstruction == nil {
		return nil, errors.New("swap transaction is empty")
	}

	//var quote QuoteResponse
	//
	//err = json.Unmarshal(result.Body(), &quote)
	//if err != nil {
	//	return nil, err
	//}

	return swap, nil
}
