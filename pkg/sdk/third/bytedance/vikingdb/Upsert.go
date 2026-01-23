package vikingdb

import (
	"context"
	"encoding/json"
	"errors"
	"store/pkg/sdk/conv"

	"github.com/go-resty/resty/v2"
)

type UpsertRequest struct {
	Collection string
	Data       map[string]interface{}
	/*
		正整数，负数无效
		当数据不过期时，默认为0。
		数据过期时间，单位为秒。设置为86400，则1天后数据自动删除。
		数据ttl删除，不会立刻更新到索引。
	*/
	TTL int
	/*
		异步写入开关

		异步写入限流阈值为同步写入的10倍
		异步写入的数据不会同步实时的写入collection，滞后时间为分钟级别。可通过接口 FetchDataInCollection来确认数据是否已经写入collection
		异步写入的数据不会触发索引的流式更新，索引同步时间为小时级别。可通过接口 FetchDataInIndex接口确认数据是否同步至index
	*/
	Async bool
}

func (t *Client) Upsert(ctx context.Context, req UpsertRequest) (*UpsertResponse, error) {

	body := map[string]interface{}{
		"collection_name": req.Collection,
		"data":            []map[string]interface{}{req.Data},
	}

	headers := sign(t.conf.AccessKeyID, t.conf.AccessKeySecret, t.conf.Service,
		t.conf.Region, t.conf.Host,
		"POST",
		"/api/vikingdb/data/upsert",
		nil,
		conv.M2B(body),
	)

	post, err := resty.New().R().SetContext(ctx).
		SetHeaders(headers).
		SetBody(body).
		Post(t.conf.BaseURL + "/api/vikingdb/data/upsert")
	if err != nil {
		return nil, err
	}

	var res Response[UpsertResponse]
	err = json.Unmarshal(post.Body(), &res)
	if err != nil {
		return nil, err
	}

	if res.Code != "Success" {
		return nil, errors.New(res.Code)
	}

	return &res.Result, nil
}

type UpsertResponse struct {
	TokenUsage map[string]Usage
}
