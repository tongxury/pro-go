package localapi

type LocalAPIResponse[T any] struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Error      string `json:"error"`
	Data       T      `json:"data"`
}

func (t *LocalAPIResponse[T]) Success() bool {
	return t.StatusCode == 0 || t.StatusCode == 200
}

type TokenMetadata struct {
	ChainId   int    `json:"chainId"`
	Address   string `json:"address"`
	ProgramId string `json:"programId"`
	LogoURI   string `json:"logoURI"`
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	Decimals  int    `json:"decimals"`
}

type ExecuteCreateOrderParams struct {
	PrivateKeyBase58 string `json:"privateKeyBase58"`
	OrderIdKeyBase58 string `json:"orderIdKeyBase58"`
	TxBase64         string `json:"txBase64"`
}

type ExecuteCreateOrderResponse struct {
	TxId string `json:"txId"`
}

type CreateLimitOrderParams struct {
	PrivateKeyBase58 string `json:"privateKeyBase58"`
	InAmountBn       string `json:"inAmountBn"`
	OutAmountBn      string `json:"outAmountBn"`
	InputMint        string `json:"inputMint"`
	OutputMint       string `json:"outputMint"`
	ExpiredAt        int64  `json:"expiredAt"`
}

type CreateLimitOrderResponse struct {
	TxId string `json:"txId"`
}

type RaydiumBuyResult struct {
	Sign         string
	AmountOut    float64
	MinAmountOut float64
	Price        float64
	//CurrentPrice   any
	//ExecutionPrice any
	//PriceImpact    any
	Fee any
}

type RaydiumSellResult struct {
	Sign         string
	AmountOut    float64
	MinAmountOut float64
	Price        float64
	//CurrentPrice   any
	//ExecutionPrice any
	//PriceImpact    any
	Fee any
}

type RaydiumBuyParams struct {
	PrivateKeyBase58 string  `json:"privateKeyBase58"`
	SolAmount        float64 `json:"solAmount"`
	TokenMint        string  `json:"tokenMint"`
	Slippage         float64 `json:"slippage"` // range: 1 ~ 0.0001, means 100% ~ 0.01%
}

type RaydiumSellParams struct {
	PrivateKeyBase58 string  `json:"privateKeyBase58"`
	Amount           float64 `json:"amount"`
	TokenMint        string  `json:"tokenMint"`
	Slippage         float64 `json:"slippage"` // range: 1 ~ 0.0001, means 100% ~ 0.01%
}

type FixedSide string

const (
	FixedSide_In  FixedSide = "in"
	FixedSide_Out FixedSide = "out"
)

type SwapToken struct {
	Address string  `json:"address"`
	Amount  float64 `json:"amount"`
}

type SwapParams struct {
	PrivateKeyBase58 string    `json:"privateKeyBase58"`
	InToken          SwapToken `json:"inToken"`
	OutToken         SwapToken `json:"outToken,omitempty"`
	FixedSide        FixedSide `json:"fixedSide,omitempty"`
	Slippage         float64   `json:"slippage,omitempty"` // range: 1 ~ 0.0001, means 100% ~ 0.01%
}

type SwapResult struct {
	Sign         string
	AmountOut    float64
	MinAmountOut float64
	Price        float64
	//CurrentPrice   any
	//ExecutionPrice any
	//PriceImpact    any
	Fee any
}
