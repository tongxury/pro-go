package volcengine

import (
	"context"
	"encoding/json"
	"errors"
)

type QueryMixCutTaskPreviewParams struct {
	TaskKey string
}

func (t *Client) QueryMixCutTaskPreview(ctx context.Context, params QueryMixCutTaskPreviewParams) (*QueryMixCutTaskPreviewResult, error) {
	bytes, err := t.doRequest(ctx, Req{
		Action: "QueryMixCutTaskPreview",
		Method: "GET",
		Params: map[string]string{
			"TaskKey":       params.TaskKey,
			"RecommendType": "1",
		},
	})
	if err != nil {
		return nil, err
	}

	var resp Response[QueryMixCutTaskPreviewResult]
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(string(bytes))
	}

	return &resp.Result, nil
}

type QueryMixCutTaskPreviewResult struct {
	Msg                string `json:"Msg"`
	RecommendTypeCount struct {
		Field1 int `json:"1"`
	} `json:"RecommendTypeCount"`
	Status  int    `json:"Status"`
	Task    []Task `json:"Task"`
	TaskKey string `json:"TaskKey"`
	Total   int    `json:"Total"`
}

type Task struct {
	CoverUrl      string  `json:"CoverUrl"`
	CreativePoint float64 `json:"CreativePoint"`
	Duration      float64 `json:"Duration"`
	GroupId       int     `json:"GroupId"`
	RecommendType string  `json:"RecommendType"`
}
