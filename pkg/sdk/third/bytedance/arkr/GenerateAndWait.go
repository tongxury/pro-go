package arkr

import (
	"context"
	"fmt"
	"store/pkg/sdk/conv"
	"time"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
)

func (t *Client) GenerateAndWait(ctx context.Context, params model.CreateContentGenerationTaskRequest) (*model.GetContentGenerationTaskResponse, error) {

	resp, err := t.c.CreateContentGenerationTask(ctx, params)
	if err != nil {
		return nil, err
	}

	for {
		time.Sleep(5 * time.Second)
		task, err := t.c.GetContentGenerationTask(ctx, model.GetContentGenerationTaskRequest{
			ID: resp.ID,
		})
		if err != nil {
			continue
		}
		fmt.Println(resp.ID, "checking task progress: "+task.Status)

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
	}

}
