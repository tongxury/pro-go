package xhs

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"time"
)

type Client struct {
	endpoint string
	lc       *cache.Cache
}

func NewClient() *Client {
	return &Client{
		endpoint: "https://adapi.xiaohongshu.com/api/open/common",
		lc:       cache.New(5*time.Minute, 10*time.Minute),
	}
}

type SendV2Params struct {
	Platform     string            `json:"platform"`
	Timestamp    int64             `json:"timestamp"`
	Scene        string            `json:"scene"`
	AdvertiserId string            `json:"advertiser_id"`
	Os           int               `json:"os"`
	Caid1Md5     string            `json:"caid1_md5"`
	ClickId      string            `json:"click_id"`
	EventType    int               `json:"event_type"`
	Context      ConversionContext `json:"context"`
}

type ConversionContext struct {
	Properties ConversionProperties `json:"properties"`
}

type ConversionProperties struct {
	ConvCnt           int   `json:"conv_cnt"`
	Pay               int   `json:"pay"`
	ActivateTimestamp int64 `json:"activate_timestamp"`
}

func (t Client) SendV2(ctx context.Context, params SendV2Params) error {

	//https://adapi.xiaohongshu.com/api/open/conversion
	re, err := resty.New().R().
		SetContext(ctx).
		SetBody(params).
		Post("https://adapi.xiaohongshu.com/api/open/conversion")

	if err != nil {
		return err
	}

	fmt.Println(re.String())

	return nil
}

// 7688267
func (t Client) Send(ctx context.Context, advertiserId, clickId, eventType string) error {

	token, err := t.getAccessToken(ctx, advertiserId)
	if err != nil {
		return err
	}

	_, err = resty.New().R().
		SetContext(ctx).
		SetBody(map[string]any{
			"advertiser_id": advertiserId,
			"method":        "aurora.leads",
			"access_token":  token,
			"event_type":    eventType,
			"conv_time":     time.Now().UnixMilli(),
			"click_id":      clickId,
		}).
		Post(t.endpoint)

	if err != nil {
		return nil
	}

	//fmt.Println(result)

	return nil
}

func (t Client) getAccessToken(ctx context.Context, advertiserId string) (string, error) {

	post, err := resty.New().R().
		SetContext(ctx).
		SetBody(map[string]string{
			"advertiser_id": advertiserId,
			"method":        "oauth.getAccessToken",
		}).
		Post(t.endpoint)

	if err != nil {
		return "", nil
	}

	var r R[struct {
		AccessToken string `json:"access_token"`
	}]

	err = json.Unmarshal(post.Body(), &r)
	if err != nil {
		return "", err
	}

	return r.Data.AccessToken, nil

}

type R[T any] struct {
	Code int    `json:"code,omitempty"`
	Msg  string `json:"msg,omitempty"`
	Data T      `json:"data,omitempty"`
}
