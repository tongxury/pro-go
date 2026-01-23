package volcengine

import (
	"context"
	"encoding/json"
	"errors"
	"store/pkg/sdk/conv"
)

type SubmitMixCutTaskAsyncParams struct {
	TaskKey  string
	GroupIds []int
}

type SubmitMixCutTaskAsyncResult struct {
	Data struct {
		TaskKey string
	}
	Message string
}

func (t *Client) SubmitMixCutTaskAsync(ctx context.Context, params SubmitMixCutTaskAsyncParams) (*SubmitMixCutTaskAsyncResult, error) {
	bytes, err := t.doRequest(ctx, Req{
		Action: "SubmitMixCutTaskAsync",
		Method: "POST",
		Params: map[string]string{
			"TaskKey": params.TaskKey,
		},
		Body: conv.M2B(map[string]any{
			"GroupIds":   params.GroupIds,
			"TaskKey":    params.TaskKey,
			"SaveRepeat": false,
		}),
	})
	if err != nil {
		return nil, err
	}

	var resp Response[SubmitMixCutTaskAsyncResult]
	err = json.Unmarshal(bytes, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 0 {
		return nil, errors.New(string(bytes))
	}

	return &resp.Result, nil
}
