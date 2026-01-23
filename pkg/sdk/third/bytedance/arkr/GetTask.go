package arkr

import (
	"context"
	"fmt"
	"store/pkg/sdk/conv"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func (t *Client) GetTask(ctx context.Context, id string) (*model.GetContentGenerationTaskResponse, error) {

	task, err := t.c.GetContentGenerationTask(ctx, model.GetContentGenerationTaskRequest{
		ID: id,
	})
	if err != nil {
		return nil, err
	}

	return &task, nil

	//queued：排队中。
	//running：任务运行中。
	//succeeded： 任务成功。（如发送失败，即5秒内没有接收到成功发送的信息，回调三次）
	//failed：任务失败。（如发送失败，即5秒内没有接收到成功发送的信息，回调三次）
	if task.Status == "succeeded" {
		return &task, nil
	}

	if task.Status == "failed" {
		return nil, fmt.Errorf("task failed: %s", conv.S2J(task))
	}

	return nil, nil
}
