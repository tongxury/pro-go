package jupiter

import (
	"context"
	"encoding/json"
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/go-resty/resty/v2"
	"net/http"
	"store/pkg/sdk/chain/sol/solana"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper/restyd"
	"time"
)

type CancelOrderRequest struct {
	Owner           string   `json:"owner"`
	FeePayer        string   `json:"feePayer"`
	OrderPublicKeys []string `json:"orders"`
}

type CancelOrderResponse struct {
	Transaction string `json:"tx"`
}

func (t *JupiterClient) CancelOrder(ctx context.Context, params CancelOrderRequest) (*CancelOrderResponse, error) {
	url := t.orderEndpoint + "/api/limit/v1/cancelorders"

	//p, _ := json.Marshal(params)

	var swap CancelOrderResponse

	result, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return nil, err
	}

	err = restyd.ParseResult(result, &swap)
	if err != nil {
		return nil, err
	}

	if swap.Transaction == "" {
		return nil, errors.New("transaction is empty")
	}

	//var quote QuoteResponse
	//
	//err = json.Unmarshal(result.Body(), &quote)
	//if err != nil {
	//	return nil, err
	//}

	return &swap, nil
}

type ListOrderParams struct {
	Wallet     string
	InputMint  string
	OutputMint string
}

type OpenOrder struct {
	PublicKey string `json:"publicKey"`
	Account   struct {
		Maker        string      `json:"maker"`
		InputMint    string      `json:"inputMint"`
		OutputMint   string      `json:"outputMint"`
		OriInAmount  string      `json:"oriInAmount"` // 原始输入数量，指的是在创建此订单时用户期望提供的代币数量
		OriOutAmount string      `json:"oriOutAmount"`
		InAmount     string      `json:"inAmount"`
		OutAmount    string      `json:"outAmount"`
		ExpiredAt    interface{} `json:"expiredAt"`
		Base         string      `json:"base"`
	} `json:"account"`
}

func (t *OpenOrder) Valid() bool {
	return t.Account.InputMint == solana.SolMint.String() || t.Account.OutputMint == solana.SolMint.String()
}

type OrderMetadata struct {
	Token          string
	TokenAmount    string
	SolAmount      string
	OriTokenAmount string
	OriSolAmount   string
	IsBuy          bool
	OrderKey       string
	ExpiredAt      int64
}

func (t *OpenOrder) Metadata() *OrderMetadata {

	if t.Account.InputMint == solana.SolMint.String() {
		return &OrderMetadata{
			Token:          t.Account.OutputMint,
			TokenAmount:    t.Account.OutAmount,
			SolAmount:      t.Account.InAmount,
			OriTokenAmount: t.Account.OriOutAmount,
			OriSolAmount:   t.Account.OriInAmount,
			IsBuy:          true,
			OrderKey:       t.PublicKey,
			ExpiredAt:      conv.Int64(t.Account.ExpiredAt),
		}
	}

	return &OrderMetadata{
		Token:          t.Account.InputMint,
		TokenAmount:    t.Account.InAmount,
		SolAmount:      t.Account.OutAmount,
		OriTokenAmount: t.Account.OriInAmount,
		OriSolAmount:   t.Account.OriOutAmount,
		IsBuy:          false,
		OrderKey:       t.PublicKey,
		ExpiredAt:      conv.Int64(t.Account.ExpiredAt),
	}

}

type Order struct {
	Id           int         `json:"id"`
	OrderKey     string      `json:"orderKey"`
	Maker        string      `json:"maker"`
	InputMint    string      `json:"inputMint"`
	OutputMint   string      `json:"outputMint"`
	InAmount     string      `json:"inAmount"`
	OriInAmount  string      `json:"oriInAmount"`
	OutAmount    string      `json:"outAmount"`
	OriOutAmount string      `json:"oriOutAmount"`
	ExpiredAt    interface{} `json:"expiredAt"`
	State        string      `json:"state"` // Cancelled  Completed
	CreateTxid   string      `json:"createTxid"`
	CancelTxid   interface{} `json:"cancelTxid"`
	UpdatedAt    time.Time   `json:"updatedAt"`
	CreatedAt    time.Time   `json:"createdAt"`
}

func (t *Order) Metadata() *OrderMetadata {

	if t.InputMint == solana.SolMint.String() {
		return &OrderMetadata{
			Token:          t.OutputMint,
			TokenAmount:    t.OutAmount,
			SolAmount:      t.InAmount,
			OriTokenAmount: t.OriOutAmount,
			OriSolAmount:   t.OriInAmount,
			IsBuy:          true,
			OrderKey:       t.OrderKey,
			ExpiredAt:      conv.Int64(t.ExpiredAt),
		}
	}

	return &OrderMetadata{
		Token:          t.InputMint,
		TokenAmount:    t.InAmount,
		SolAmount:      t.OutAmount,
		OriTokenAmount: t.OriInAmount,
		OriSolAmount:   t.OriOutAmount,
		IsBuy:          false,
		OrderKey:       t.OrderKey,
		ExpiredAt:      conv.Int64(t.ExpiredAt),
	}

}

type OpenOrders []*OpenOrder

func (ts OpenOrders) AsMap() map[string]*OpenOrder {

	result := make(map[string]*OpenOrder, len(ts))
	for i := range ts {
		x := ts[i]

		result[x.PublicKey] = x
	}

	return result
}

func (ts OpenOrders) TokenIds() []string {

	var tokenIds []string
	for _, x := range ts {
		tokenIds = append(tokenIds, x.Account.InputMint, x.Account.OutputMint)
	}

	return mapset.NewSet(tokenIds...).ToSlice()
}

type Orders []*Order

func (ts Orders) TokenIds() []string {

	var tokenIds []string
	for _, x := range ts {
		tokenIds = append(tokenIds, x.InputMint, x.OutputMint)
	}

	return mapset.NewSet(tokenIds...).ToSlice()
}

type ListOrdersParams struct {
	Wallet string
	Cursor string
	Skip   string
	Take   int
}

func (t *JupiterClient) ListOrders(ctx context.Context, params ListOrdersParams) (Orders, error) {

	url := t.orderEndpoint + "/api/limit/v1/orderHistory"

	var orders Orders

	rsp, err := resty.New().R().SetQueryParams(map[string]string{
		"wallet": params.Wallet,
		"cursor": params.Cursor,
		"skip":   params.Skip,
		"take":   conv.String(params.Take),
	}).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(rsp.Body()))
	}

	err = json.Unmarshal(rsp.Body(), &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (t *JupiterClient) ListOpenOrders(ctx context.Context, params ListOrderParams) (OpenOrders, error) {

	url := t.orderEndpoint + "/api/limit/v1/openorders"

	var orders OpenOrders

	rsp, err := resty.New().R().SetQueryParams(map[string]string{
		"inputMint":  params.InputMint,
		"outputMint": params.OutputMint,
		"wallet":     params.Wallet,
	}).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	if rsp.StatusCode() != http.StatusOK {
		return nil, errors.New(string(rsp.Body()))
	}

	err = json.Unmarshal(rsp.Body(), &orders)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

type CreateOrderRequest struct {
	Owner           string `json:"owner"`
	InAmount        string `json:"inAmount"`
	OutAmount       string `json:"outAmount"`
	InputMint       string `json:"inputMint"`
	OutputMint      string `json:"outputMint"`
	ExpiredAtSecond int64  `json:"expiredAt"` // unix
	OrderId         string `json:"base"`
}

type CreateOrderResponse struct {
	Transaction string `json:"tx"`
	OrderPubkey string `json:"orderPubkey"`
}

func (t *JupiterClient) CreateOrder(ctx context.Context, params CreateOrderRequest) (*CreateOrderResponse, error) {
	url := t.orderEndpoint + "/api/limit/v1/createOrder"

	//p, _ := json.Marshal(params)

	var swap CreateOrderResponse

	result, err := resty.New().R().
		SetHeader("Content-Type", "application/json").
		SetBody(params).
		SetContext(ctx).
		Post(url)

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(result.Body(), &swap)

	if err != nil {
		return nil, err
	}

	if swap.Transaction == "" {
		return nil, errors.New("transaction is empty")
	}

	//var quote QuoteResponse
	//
	//err = json.Unmarshal(result.Body(), &quote)
	//if err != nil {
	//	return nil, err
	//}

	return &swap, nil
}
