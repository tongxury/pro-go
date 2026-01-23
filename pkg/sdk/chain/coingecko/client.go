package coingecko

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"store/pkg/sdk/helper/restyd"
	"strings"
)

var (
	CoinNotFound = fmt.Errorf("coin not found")
)

type Client struct {
	endpoint string
}

func NewClient() *Client {
	return &Client{
		endpoint: "https://api.coingecko.com",
	}
}

func (c *Client) GetCoinId(ctx context.Context, address string) (string, error) {
	return "", nil
}

type ErrorResult struct {
	Error string `json:"error"`
}

func (c *Client) GetPrices(ctx context.Context, address string, tsFrom, tsTo int64) (Prices, error) {

	url := fmt.Sprintf("%s/api/v3/coins/ethereum/contract/%s/market_chart/range?vs_currency=usd&from=%d&to=%d",
		c.endpoint, address, tsFrom, tsTo,
	)

	response, err := resty.New().R().Get(url)

	if err != nil {
		if strings.Contains(err.Error(), "coin not found") {
			return nil, CoinNotFound
		}
		return nil, err
	}

	var tmp struct {
		Prices     [][]float64
		MarketCaps [][]float64
	}

	err = restyd.ParseResult(response, &tmp)
	if err != nil {
		return nil, err
	}

	var rsp Prices
	for _, x := range tmp.Prices {
		rsp = append(rsp, &Price{
			Address: address,
			Ts:      int64(x[0] / 1000),
			Price:   x[1],
		})
	}

	return rsp, nil
}

func (c *Client) GetCoinInfo(ctx context.Context, address string) (string, error) {
	return "", nil
}

//func (c *Client) ListCoins(ctx context.Context) (Coins, error) {
//	url := fmt.Sprintf("%s/api/v3/coins/list", c.base)
//
//	bytes, code, err := httpcli.Client().GET(ctx, url)
//	if err != nil {
//		return nil, err
//	}
//	if code != 200 {
//		return nil, fmt.Errorf(string(bytes))
//	}
//
//	var rsp Coins
//	err = conv.B2S(bytes, &rsp)
//	if err != nil {
//		return nil, err
//	}
//
//	return rsp, nil
//}

func (c *Client) GetCoinOhlcs(ctx context.Context, id string, days int64) (Ohlcs, error) {
	url := fmt.Sprintf("%s/api/v3/coins/%s/ohlc?vs_currency=usd&days=%d", c.endpoint, id, days)

	response, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}

	var parts [][]float64

	err = restyd.ParseResult(response, &parts)
	if err != nil {
		return nil, err
	}
	var rsp Ohlcs
	for _, part := range parts {
		rsp = append(rsp, &Ohlc{
			Id:    id,
			Time:  int64(part[0] / 1000),
			Open:  part[1],
			High:  part[2],
			Low:   part[3],
			Close: part[4],
		})
	}

	return rsp, nil
}
