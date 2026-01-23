package gmgnai

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
)

func (t *Client) GetLatestTrades(ctx context.Context, params GetLatestTradesParams) (*LatestTrades, error) {

	url := fmt.Sprintf("https://gmgn.ai/defi/quotation/v1/trades/%s/%s", params.Chain, params.Token)

	// todo  每次cookie都不一样
	//cookie := "_ga=GA1.1.1541192268.1717413851; __cf_bm=tyGnynE6u_s7qRuq.mkBhYczw9Muo.KTrCwapbLRdV4-1731122206-1.0.1.1-zmeKCbadCvGTsD8aLKaCUiEWjRCQ2U1NIwjKz4GRV7wyDMqYrMzvlytU8rbmzbz5uEXHPDzLDcKMOLeJ.QmZrQ; cf_clearance=l_gZBOPnjHTDMMeFtRm9ItfwB00n1Mm5y6XXBgDHjw8-1731122400-1.2.1.1-qD2RY3JMcRW180ojesSh8h5JCYhjWg1cMr87gn52mTzSGAysrVFBKEaF1xFB.DhHKXqvx9WjQrRzmpz.KwWX_TvDj9YAKKkOkxi7mdgLDgjotW3wUBeOBTaxufyeqmOsy1m0ZPGNv_bxYnnnzSY7fCWLwZ7sjrGvYVt3zytnRe3Lump47thr.ocBl1e8uLNN0jnbSGol5zbG7hmAH_.JxustWhNg6.HP_bf9CwqW1RhYaihpeUbVz4M_hDtypE7jA8Z1NxJQ4G5XZpzuFCMwCEIw3VneSXLB3qWKTGpf.OiZKtZ2Um7mYWjBcVVbu3RvBWRmrd3Dzg_8aV3zA72qcf9A2yQjB79_eVINYq5t1lowo6h6ql9Aa0aRGVwjnot_CFvc1b.JXuuL_i1v1Z4G7w; _ga_0XM0LYXGC8=GS1.1.1731121305.342.1.1731122560.0.0.0"
	//cookie := "_ga=GA1.1.1541192268.1717413851; __cf_bm=pCI2gKfGWW8CvF5ATZogCe7Cp.B_.9DoNeBpw4tXZyM-1731124993-1.0.1.1-w_rZQPjcu97gd1KqiGdaAvvAdr4BsRXY5vdm2pDbv4jcL.GG71rOSS9_bf_2RsQCLhKgEISiHJHb0OR.XqgH9Q; cf_clearance=p0IZ7_MEh2HQhLdROYnFJ1kVOK.0IlqDCoGuEP38uRw-1731125744-1.2.1.1-.QgARc4n_a.OByKrHBNdeKG46sfIzItVDq4l3UfWcap3QjnxifFIkvUJEmVmIS1bqYov.4kIfwoJlQKpbprafg9pmy5xQtyFsUuvkUj.kyr49vcyhtvnXYfBhbA9ZuXNxW93xtOB_OzLEy2udzcmMMpjIBTHtY1C1bLK_LrsyRs7i5pHZr.6bKsaFLDOf_K1rWonUYIfbIX4716C2cqAOe4iHRkYjuWQzpdZIuYXXRXdI_4ZQH99rzO2PZ7nR7XnXJqA3UxDfb6UN6Ln1FCWcCaqIZdUn6.CRvnynBcS8TVvOtNp9emrcf61bKB_Y3k.i8LOf485pmtTgjeZZNHqajw9whcu6CBP8eST77P6qRuSRCOa_Xa576ZTW_ZY02I0BZrEYWgFwt5qRMNuNp1Buw; _ga_0XM0LYXGC8=GS1.1.1731125736.343.1.1731125749.0.0.0"
	result, err := resty.New().R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"limit": fmt.Sprintf("%d", params.Limit),
			"maker": params.Maker,
		}).
		SetHeaders(map[string]string{
			//"Cookie":             cookie,
			"Content-Type":       "application/json",
			"sec-ch-ua":          "\"Chromium\";v=\"130\", \"Google Chrome\";v=\"130\", \"Not?A_Brand\";v=\"99\"",
			"sec-ch-ua-platform": "\"macOS\"",
			"sec-ch-ua-mobile":   "?0",
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "same-origin",
			"priority":           "u=1, i",
			"pragma":             "no-cache",
			"Accept":             "application/json, text/plain, */*",
			"Accept-Language":    "zh-CN,zh;q=0.9,ar;q=0.8",
			"Referer":            "https://gmgn.ai/sol/token/" + params.Token,
			"User-Agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/130.0.0.0 Safari/537.36",
		}).
		Get(url)
	if err != nil {
		return nil, err
	}

	if result.IsError() {
		return nil, errors.New(result.String())
	}

	var trades Response[LatestTrades]
	err = json.Unmarshal(result.Body(), &trades)
	if err != nil {
		return nil, err
	}

	if trades.Code > 0 {
		return nil, errors.New(trades.Msg)
	}

	return &trades.Data, nil
}

type GetLatestTradesParams struct {
	Chain string
	Token string
	Maker string
	Limit int
}

type Trade struct {
	Maker                string      `json:"maker"`
	BaseAmount           float64     `json:"base_amount"`
	QuoteAmount          float64     `json:"quote_amount"`
	AmountUsd            *float64    `json:"amount_usd"`
	Timestamp            int         `json:"timestamp"`
	Event                string      `json:"event"`
	TxHash               string      `json:"tx_hash"`
	PriceUsd             *string     `json:"price_usd"`
	MakerTags            []string    `json:"maker_tags"`
	MakerTwitterUsername interface{} `json:"maker_twitter_username"`
	MakerTwitterName     interface{} `json:"maker_twitter_name"`
	MakerName            interface{} `json:"maker_name"`
	MakerAvatar          interface{} `json:"maker_avatar"`
	MakerEns             interface{} `json:"maker_ens"`
	MakerTokenTags       []string    `json:"maker_token_tags"`
	TokenAddress         string      `json:"token_address"`
	QuoteAddress         string      `json:"quote_address"`
	TotalTrade           int         `json:"total_trade"`
	Id                   string      `json:"id"`
	IsFollowing          int         `json:"is_following"`
	IsOpenOrClose        int         `json:"is_open_or_close"`
	BuyCostUsd           *float64    `json:"buy_cost_usd"`
	Balance              string      `json:"balance"`
	Cost                 float64     `json:"cost"`
	HistoryBoughtAmount  float64     `json:"history_bought_amount"`
	HistorySoldIncome    float64     `json:"history_sold_income"`
	HistorySoldAmount    float64     `json:"history_sold_amount"`
	UnrealizedProfit     float64     `json:"unrealized_profit"`
	RealizedProfit       float64     `json:"realized_profit"`
	TokenSymbol          string      `json:"token_symbol,omitempty"`
	QuoteSymbol          string      `json:"quote_symbol,omitempty"`
}

type LatestTrades struct {
	History []Trade `json:"history"`
	Next    string  `json:"next"`
}
