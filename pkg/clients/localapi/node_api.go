package localapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	endpoint string
}

func NewClient(endpoint string) *Client {

	return &Client{
		endpoint: endpoint,
	}
}

type CreateTradeTransactionParams struct {
	Token             string  `json:"token"`
	Side              string  `json:"side"`
	AmountWithDecimal int64   `json:"amount"`
	Slippage          float64 `json:"slippage"`
	PriorityFee       float64 `json:"priorityFee"`
	Wallet            string  `json:"wallet"`
	PrivateKey        string  `json:"privateKey"`
}

type CreateTradeTransactionResponse struct {
	RawTx string `json:"rawTx"`
}

func (t *Client) CreateTradeTransaction(ctx context.Context, params CreateTradeTransactionParams) (*CreateTradeTransactionResponse, error) {
	result, err := resty.New().R().SetContext(ctx).
		SetBody(params).
		Post(t.endpoint + "/api/v1/trade")

	if err != nil {
		return nil, err
	}

	var response LocalAPIResponse[CreateTradeTransactionResponse]
	err = json.Unmarshal(result.Body(), &response)
	if err != nil {
		return nil, err
	}

	if !response.Success() {
		return nil, errors.New(result.String())
	}

	return &response.Data, nil
}
