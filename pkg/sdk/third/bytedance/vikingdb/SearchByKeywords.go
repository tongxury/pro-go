package vikingdb

import (
	"context"
	"encoding/json"
	"store/pkg/sdk/conv"

	"github.com/go-resty/resty/v2"
)

// https://www.volcengine.com/docs/84313/1791133
type SearchRequest struct {
	CollectionName string         `json:"collection_name,omitempty"`
	IndexName      string         `json:"index_name,omitempty"`
	OutputFields   []string       `json:"output_fields,omitempty"`
	Filter         map[string]any `json:"filter,omitempty"`
	Limit          int            `json:"limit,omitempty"`
	Offset         int            `json:"offset,omitempty"`
}

type SearchByKeywordsRequest struct {
	SearchRequest
	Keywords      []string `json:"keywords,omitempty"`
	CaseSensitive bool     `json:"case_sensitive,omitempty"`
}

func (t *Client) SearchByKeywords(ctx context.Context, req SearchByKeywordsRequest) (*SearchResponse, error) {

	body, _ := conv.S2M[string, any](req)

	path := "/api/vikingdb/data/search/keywords"

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

	return &res.Result, nil
}

type Item struct {
	Id     string `json:"id"`
	Fields struct {
		FText string `json:"f_text"`
	} `json:"fields"`
	Score    float64 `json:"score"`
	AnnScore float64 `json:"ann_score"`
}

type SearchResponse struct {
	Data             []Item
	TotalReturnCount int
	TokenUsage       map[string]Usage
}
