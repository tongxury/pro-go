package gmgnai

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"store/pkg/sdk/helper/restyd"
	"strconv"
	"time"
)

type OHLC struct {
	Open   string `json:"open"`
	High   string `json:"high"`
	Low    string `json:"low"`
	Close  string `json:"close"`
	Volume string `json:"volume"`
	Time   string `json:"time"`
}

type OHLCs []*OHLC

type GetOHLCsParams struct {
	Token      string
	Resolution string // 1d 1m
	FromTs     int64
	ToTs       int64
}

func (t *Client) GetPrice(ctx context.Context, token string) (float64, error) {
	cs, err := t.GetOHLCs(ctx, GetOHLCsParams{
		Token:      token,
		Resolution: "1m",
		FromTs:     time.Now().Add(10 * time.Second).Unix(),
		ToTs:       time.Now().Unix(),
	})
	if err != nil {
		return 0, err
	}

	if len(cs) == 0 {
		return 0, errors.New("no price found: " + token)
	}

	result, _ := strconv.ParseFloat(cs[len(cs)-1].Close, 64)

	return result, nil
}

func (t *Client) GetOHLCs(ctx context.Context, params GetOHLCsParams) (OHLCs, error) {

	url := fmt.Sprintf("https://www.gmgn.cc/defi/quotation/v1/tokens/kline/sol/%s?resolution=%s&from=%d&to=%d",
		params.Token, params.Resolution, params.FromTs, params.ToTs)

	response, err := resty.New().R().SetContext(ctx).Get(url)
	if err != nil {
		return nil, err
	}

	var data Response[OHLCs]
	err = restyd.ParseResult(response, &data)
	if err != nil {
		return nil, err
	}

	return data.Data, nil
}
