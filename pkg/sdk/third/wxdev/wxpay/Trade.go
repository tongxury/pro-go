package wxpay

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
)

type TradeParams struct {
	Description string
	OutTradeNo  string
	NotifyUrl   string
	TotalAmount int64 // 分
	OpenId      string
}

type Transaction struct {
	PrepayId  string `json:"prepay_id"`
	Timestamp int64
	Nonce     string
	Package   string
	SignType  string
	PaySign   string
}

func (t *Client) Trade(ctx context.Context, params TradeParams) (*Transaction, error) {

	//var mchID string = "1725227465"
	//var mchCertificateSerialNumber string = "1F523651CF31057ACACA4E5DEBD292F5081D6C15"
	////var pubKeyId = "PUB_KEY_ID_0117252274652025081700192211000800"
	//var mchAPIv3Key string = "MIIEvQIBADANBgkqhkiG9w0BAQEFAASA"

	appId := "wxb0008ff0fb65a736"

	p := map[string]any{
		"appid":        appId,
		"mchid":        "1725227465",
		"description":  params.Description,
		"out_trade_no": params.OutTradeNo,
		"notify_url":   params.NotifyUrl,
		"amount": map[string]any{
			"total":    params.TotalAmount,
			"currency": "CNY",
		},
		"payer": map[string]any{
			"openid": params.OpenId,
		},
	}

	//appid = "wxd678efh567hg6787"
	//mchid = "1900007291"
	//description = "Image形象店-深圳腾大-QQ公仔"
	//out_trade_no = "1217752501201407033233368018"
	//notify_url = "https://www.weixin.qq.com/wxpay/pay.php"
	//amount = {"total":100,"currency":"CNY"}
	//payer = {"openid":"oUpF8uMuAJO_M2pxb1Q9zNjWeS6o"}

	body, _ := json.Marshal(p)

	bd := string(body)

	signResult, err := t.sign(
		"POST", "/v3/pay/transactions/jsapi", bd)
	if err != nil {
		return nil, err
	}

	post, err := resty.New().R().
		SetContext(ctx).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", signResult.Signature).
		SetBody(body).
		Post("https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi")

	if err != nil {
		return nil, err
	}

	fmt.Println("result", string(post.Body()))

	var transaction Transaction
	err = json.Unmarshal(post.Body(), &transaction)
	if err != nil {
		return nil, err
	}

	pkg := "prepay_id=" + transaction.PrepayId

	signMessage := fmt.Sprintf(`"%s
%d
%s
%s
"`, appId, signResult.TimeStamp, signResult.NonceStr, pkg)

	signMessage = fmt.Sprintf("%s\n%d\n%s\n%s\n", appId, signResult.TimeStamp, signResult.NonceStr, pkg)

	message, err := t.signMessage(signMessage)
	if err != nil {
		return nil, err
	}

	fmt.Println("signMessage", signMessage)

	return &Transaction{
		PrepayId:  transaction.PrepayId,
		Timestamp: signResult.TimeStamp,
		Nonce:     signResult.NonceStr,
		Package:   pkg,
		SignType:  "RSA",
		PaySign:   message,
	}, nil
}
