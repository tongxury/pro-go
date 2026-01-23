package quicknode

import (
	"context"
	"encoding/json"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	endpoint string
}

func NewClient() *Client {
	return &Client{
		endpoint: "https://light-proud-pine.solana-mainnet.quiknode.pro/ab05da0bef752cdf59801f675a549691dc45e4c6",
	}
}

func (t *Client) EstimatePriorityFees(ctx context.Context, params *EstimatePriorityFeesParams) (*PriorityFees, error) {

	params.ApiVersion = 2
	if params.LastNBlocks == 0 {
		params.LastNBlocks = 100
	}

	result, err := resty.New().R().
		SetBody(Params{
			Jsonrpc: "2.0",
			Id:      1,
			Method:  "qn_estimatePriorityFees",
			Params:  *params,
		}).
		SetContext(ctx).
		Post(t.endpoint)

	if err != nil {
		return nil, err
	}

	var rsp Result
	err = json.Unmarshal(result.Body(), &rsp)
	if err != nil {
		return nil, err
	}

	return &rsp.Result, nil
}

type EstimatePriorityFeesParams struct {
	LastNBlocks int    `json:"last_n_blocks"`
	Account     string `json:"account"`
	ApiVersion  int    `json:"api_version"`
}

type Params struct {
	Jsonrpc string                     `json:"jsonrpc"`
	Id      int                        `json:"id"`
	Method  string                     `json:"method"`
	Params  EstimatePriorityFeesParams `json:"params"`
}

type PriorityFees struct {
	Context struct {
		Slot int `json:"slot"`
	} `json:"context"`
	PerComputeUnit T `json:"per_compute_unit"`
	PerTransaction T `json:"per_transaction"`
}

type Result struct {
	Jsonrpc string       `json:"jsonrpc"`
	Result  PriorityFees `json:"result"`
	Id      int          `json:"id"`
}

type T struct {
	Extreme     int64         `json:"extreme"`
	High        int64         `json:"high"`
	Low         int64         `json:"low"`
	Medium      int64         `json:"medium"`
	Percentiles map[int]int64 `json:"percentiles"`
}
