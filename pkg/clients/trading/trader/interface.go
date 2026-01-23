package trader

import (
	"context"
	"math"
	"store/pkg/sdk/chain/sol/solana"
)

type ITrader interface {
	CreateInstructions(ctx context.Context, params *CreateCreateInstructionsParams) (*CreateInstructionsResult, error)
	CreateTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateTransactionResult, error)
	CreateRawTransaction(ctx context.Context, params *CreateTransactionParams) (*CreateRawTransactionResult, error)
}

type CreateRawTransactionResult struct {
	RawTX *solana.Transaction `json:"tx"`
	//Quotation Quotation           `json:"quotation"`
}
type CreateInstructionsResult struct {
	Instructions []solana.Instruction
	Quotation    *Quotation
}

type CreateCreateInstructionsParams struct {
	Token    string
	Side     string // buy or sell
	Amount   int64
	Slippage float64   // eg. 0.1 代表 10%
	Priority *Priority // eg. 0.001 sol
	Wallet   string
}

func (t *CreateCreateInstructionsParams) InputMint() string {
	if t.IsBuy() {
		return solana.SolMintString
	}

	return t.Token
}

func (t *CreateCreateInstructionsParams) OutputMint() string {
	if t.IsBuy() {
		return t.Token
	}

	return solana.SolMint.String()
}

func (t *CreateCreateInstructionsParams) IsBuy() bool {
	return t.Side == "buy"
}

// "feeBps" 通常指的是“费用基点”（fee basis points），它是在金融和投资领域用来表示费用或佣金的单位。1个基点等于0.01%（或1/100的百分比）。
type CreateTransactionParams struct {
	Token      string
	Side       string // buy or sell
	Amount     int64
	Slippage   float64   // eg. 0.1 代表 10%
	Priority   *Priority // eg. 0.001 sol
	Wallet     string
	PrivateKey string
	Fee        *Fee
}

type Priority struct {
	Fee   float64
	Level string
}

type Fee struct {
	Rate     float64
	Receiver string
}

func (t *CreateTransactionParams) InputMint() string {
	if t.IsBuy() {
		return solana.SolMintString
	}

	return t.Token
}

func (t *CreateTransactionParams) OutputMint() string {
	if t.IsBuy() {
		return t.Token
	}

	return solana.SolMint.String()
}

func (t *CreateTransactionParams) IsBuy() bool {
	return t.Side == "buy"
}

type CreateTransactionResult struct {
	TxHash    string
	Confirmed bool
	//Quote     *Quote
	Quotation *Quotation
}

type Quotation struct {
	Token       string
	Amount      int64
	QuoteAmount int64
	Side        string
}

//type Quote struct {
//	Token          string
//	TokenAmountRaw string
//	SolAmountRaw   string
//	Side           string
//}

func (t *Quotation) Price(decimals int64) float64 {
	return t.SolAmount() / t.TokenAmount(decimals)
}

func (t *Quotation) SolAmount() float64 {
	return float64(t.QuoteAmount) / math.Pow10(9)
}

func (t *Quotation) TokenAmount(decimals int64) float64 {
	return float64(t.Amount) / math.Pow10(int(decimals))
}
