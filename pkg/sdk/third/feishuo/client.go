package feishuo

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/patrickmn/go-cache"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"time"
)

type Client struct {
	endpoint string
	c        *cache.Cache
}

func NewClient(endpoint string) *Client {
	return &Client{
		endpoint: endpoint,
		c:        cache.New(time.Minute, time.Minute),
	}
}

type RichText struct {
	Key     string              `json:"key"` // 用于频率限制
	Title   string              `json:"title"`
	Content []RichTextParagraph `json:"content"`
}

type RichTextParagraph []RichTextContent

type RichTextContent struct {
	Tag    string `json:"tag"`
	Text   string `json:"text"`
	Href   string `json:"href"`
	UserID string `json:"user_id"`
}

type richTextParams struct {
	MsgType string        `json:"msg_type"`
	Content contentParams `json:"content"`
}

type contentParams struct {
	Post contentPost `json:"post"`
}

type contentPost struct {
	ZhCN RichText `json:"zh_cn"`
}

type Option struct {
	IntervalSeconds int
}

func (t *Client) SendRichText(ctx context.Context, params RichText, options ...Option) error {
	if len(options) > 0 {
		if options[0].IntervalSeconds > 0 {

			cacheKey := helper.OrString(params.Key, "key")

			_, found := t.c.Get(cacheKey)
			if found {
				return nil
			}

			err := t.sendRichText(ctx, params)
			if err != nil {
				return err
			}

			t.c.Set(cacheKey, "1", time.Duration(options[0].IntervalSeconds)*time.Second)
			return nil
		}
	}

	return t.sendRichText(ctx, params)
}

func (t *Client) sendRichText(ctx context.Context, params RichText) error {

	//log.Debugw("sendRichText", params)

	p := conv.S2J(richTextParams{
		MsgType: "post",
		Content: contentParams{
			Post: contentPost{
				ZhCN: params,
			},
		},
	})

	rsp, err := resty.New().R().SetContext(ctx).
		SetBody(p).
		Post(t.endpoint)
	if err != nil {
		return err
	}

	var result *Result[any]
	_ = conv.J2S(rsp.Body(), &result)

	if result.Code != 0 {
		return fmt.Errorf(result.Msg)
	}

	return nil
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *T     `json:"data"`
}
