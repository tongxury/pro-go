package jupiter

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/helper/restyd"
	"strings"
)

type GetPricesParams struct {
	PublicKeys []string
	//VsToken    string
	ShowExtraInfo bool
}

type Price struct {
	Id    string `json:"id"`
	Type  string `json:"type"`
	Price string `json:"price"`
}

type GetPricesResponse struct {
	Data      map[string]*Price `json:"data"`
	TimeTaken float64           `json:"timeTaken"`
}

func (t *JupiterClient) GetPrices(ctx context.Context, params GetPricesParams) (*GetPricesResponse, error) {

	url := t.apiEndpoint + "/price/v2"

	var quote GetPricesResponse

	rsp, err := resty.New().R().SetQueryParams(map[string]string{
		"ids":           strings.Join(params.PublicKeys, ","),
		"showExtraInfo": helper.Select(params.ShowExtraInfo, "true", "false"),
	}).
		SetContext(ctx).
		Get(url)

	if err != nil {
		return nil, err
	}

	err = restyd.ParseResult(rsp, &quote)
	if err != nil {
		return nil, err
	}

	if len(quote.Data) != len(params.PublicKeys) {
		return nil, fmt.Errorf("invalid result of public keys %v %v", strings.Join(params.PublicKeys, ","), quote.Data)
	}

	return &quote, nil
}
