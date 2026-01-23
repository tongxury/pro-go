package vikingdb

import (
	"context"
	"encoding/json"
	"store/pkg/sdk/conv"

	"github.com/go-resty/resty/v2"
)

// https://www.volcengine.com/docs/84313/1902648?lang=zh

type SearchByMultiModalRequest struct {
	SearchRequest
	Text            string `json:"text"`
	NeedInstruction bool   `json:"need_instruction"`
}

func (t *Client) SearchByMultiModal(ctx context.Context, req SearchByMultiModalRequest) (*SearchResponse, error) {

	body, _ := conv.S2M[string, any](req)

	path := "/api/vikingdb/data/search/multi_modal"

	headers := sign(t.conf.AccessKeyID, t.conf.AccessKeySecret, t.conf.Service,
		t.conf.Region, t.conf.Host,
		"POST",
		path,
		nil,
		conv.M2B(body),
	)

	post, err := resty.New().R().SetContext(ctx).
		SetHeaders(headers).
		SetBody(body).
		Post(t.conf.BaseURL + path)
	if err != nil {
		return nil, err
	}

	var res Response[SearchResponse]
	err = json.Unmarshal(post.Body(), &res)
	if err != nil {
		return nil, err
	}

	//if res.Code

	return &res.Result, nil
}
