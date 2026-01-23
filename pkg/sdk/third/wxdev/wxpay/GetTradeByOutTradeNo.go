package wxpay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
)

func (t *Client) GetTradeByOutTradeNo(ctx context.Context, no string) (*Trade, error) {

	sign, err := t.sign("GET", fmt.Sprintf("/v3/pay/transactions/out-trade-no/%s?mchid=%s", no, t.conf.MchId), "")
	if err != nil {
		return nil, err
	}

	post, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", sign.Signature).
		SetQueryParams(map[string]string{
			"mchid": t.conf.MchId,
		}).
		Get("https://api.mch.weixin.qq.com/v3/pay/transactions/out-trade-no/" + no)

	if err != nil {
		return nil, err
	}

	fmt.Println(string(post.Body()))
	var resp Trade
	err = json.Unmarshal(post.Body(), &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

type Trade struct {
	Amount struct {
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
		PayerTotal    int    `json:"payer_total"`
		Total         int    `json:"total"`
	} `json:"amount"`
	Appid      string `json:"appid"`
	Attach     string `json:"attach"`
	BankType   string `json:"bank_type"`
	Mchid      string `json:"mchid"`
	OutTradeNo string `json:"out_trade_no"`
	Payer      struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	PromotionDetail []interface{} `json:"promotion_detail"`
	SuccessTime     time.Time     `json:"success_time"`
	TradeState      string        `json:"trade_state"`
	TradeStateDesc  string        `json:"trade_state_desc"`
	TradeType       string        `json:"trade_type"`
	TransactionId   string        `json:"transaction_id"`
}
