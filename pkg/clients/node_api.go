package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

type NodeApi struct {
	endpoint string
}

func NewNodeApi(endpoint string) *NodeApi {

	return &NodeApi{
		endpoint: endpoint,
	}
}

func (t *NodeApi) CreateLimitOrders(ctx context.Context, params CreateLimitOrderParams) (*CreateLimitOrderResponse, error) {

	response, err := resty.New().R().SetContext(ctx).
		SetBody(params).
		Post(t.endpoint + "/api/v1/limit/createOrders")

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, fmt.Errorf(response.String())
	}

	var resp NodeApiResponse[CreateLimitOrderResponse]
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf(resp.Message)
	}

	return &resp.Data, nil
}

func (t *NodeApi) ExecuteCreateOrderTransaction(ctx context.Context, params ExecuteCreateOrderParams) (*ExecuteCreateOrderResponse, error) {

	response, err := resty.New().
		SetTimeout(120 * time.Second).
		R().
		SetContext(ctx).
		SetBody(params).
		Post(t.endpoint + "/api/v1/limit/executeOrders")

	if err != nil {
		return nil, err
	}

	if !response.IsSuccess() {
		return nil, fmt.Errorf(response.String())
	}

	var resp NodeApiResponse[ExecuteCreateOrderResponse]
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf(resp.Message)
	}

	return &resp.Data, nil
}

func (t *NodeApi) Swap(ctx context.Context, params SwapParams) (*SwapResult, int, error) {

	response, err := resty.New().R().SetContext(ctx).
		//SetQueryParam("address", addressBase58).
		SetBody(params).
		Post(t.endpoint + "/api/swaps")

	if err != nil {
		return nil, 0, err
	}

	if !response.IsSuccess() {
		return nil, response.StatusCode(), fmt.Errorf(response.String())
	}

	var resp NodeApiResponse[SwapResult]
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return nil, 0, err
	}

	if !resp.Success() {
		return nil, 0, fmt.Errorf(resp.Message)
	}

	return &resp.Data, 0, nil
}

func (t *NodeApi) FindTokenMetadata(ctx context.Context, publicKey string) (*TokenMetadata, error) {

	response, err := resty.New().R().SetContext(ctx).
		//SetQueryParam("address", addressBase58).
		Get(t.endpoint + "/api/tokens/" + publicKey + "/metadata")

	if err != nil {
		return nil, err
	}

	var resp NodeApiResponse[TokenMetadata]
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return nil, err
	}

	if !resp.Success() {
		return nil, fmt.Errorf(response.String())
	}

	return &resp.Data, nil
}

func (t *NodeApi) GetTokenPrice(ctx context.Context, publicKey string) (float64, error) {

	response, err := resty.New().R().SetContext(ctx).
		//SetQueryParam("address", addressBase58).
		Get(t.endpoint + "/api/tokens/" + publicKey + "/prices")

	if err != nil {
		return 0, err
	}

	var resp NodeApiResponse[float64]
	err = json.Unmarshal(response.Body(), &resp)
	if err != nil {
		return 0, err
	}

	if !resp.Success() {
		return 0, fmt.Errorf(response.String())
	}

	return resp.Data, nil
}
