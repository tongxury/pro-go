package volcengine

import (
	"context"
	"encoding/json"
	"errors"
)

type QueryMixCutTaskResultParams struct {
	TaskKey string
}

type QueryMixCutTaskResultResult struct {
	Data struct {
		Task struct {
			CreatedAt  string `json:"CreatedAt"`
			Status     int    `json:"Status"`
			TaskKey    string `json:"TaskKey"`
			UpdatedAt  string `json:"UpdatedAt"`
			VideoCount int    `json:"VideoCount"`
			VideoKey   string `json:"VideoKey"`
			VideoList  []struct {
				CoverUrl    string  `json:"CoverUrl"`
				CreatedAt   string  `json:"CreatedAt"`
				DownloadUrl string  `json:"DownloadUrl"`
				Duration    float64 `json:"Duration"`
				MuseId      string  `json:"MuseId"`
				TaskStatus  int     `json:"TaskStatus"`
				UpdatedAt   string  `json:"UpdatedAt"`
				VideoKey    string  `json:"VideoKey"`
			} `json:"VideoList"`
		}
	}
}

func (t *Client) QueryMixCutTaskResult(ctx context.Context, params QueryMixCutTaskResultParams) (*QueryMixCutTaskResultResult, error) {
	bytes, err := t.doRequest(ctx, Req{
		Action: "QueryMixCutTaskResult",
		Method: "GET",
		Params: map[string]string{
			"TaskKey": params.TaskKey,
		},
	})
	if err != nil {
		return nil, err
	}

	var resp Response[QueryMixCutTaskResultResult]
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(resp.Message)
	}

	return &resp.Result, nil
}
