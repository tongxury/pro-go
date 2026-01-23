package arkr

import (
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func (t *Client) GetContentGenerationTask(ctx context.Context, id string) (*model.GetContentGenerationTaskResponse, error) {

	req := model.GetContentGenerationTaskRequest{
		ID: id,
	}
	resp, err := t.c.GetContentGenerationTask(ctx, req)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}
