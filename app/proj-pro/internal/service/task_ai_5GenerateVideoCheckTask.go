package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper/wg"
	"store/pkg/sdk/third/bytedance/volcengine"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) GenerateVideoCheckTask() {

	ctx := context.Background()

	list, err := t.data.Mongo.Task.List(ctx, bson.M{"status": "videoGeneratingWaiting"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.generateVideoCheckTask)
}

func (t ProjService) generateVideoCheckTask(ctx context.Context, task *projpb.Task) error {

	log.Debugw("generateVideoCheckTask", task.XId)

	result, err := t.data.Volcengine.QueryMixCutTaskResult(ctx, volcengine.QueryMixCutTaskResultParams{
		TaskKey: task.GeneratedResult.GetTaskId(),
	})
	if err != nil {
		return err
	}

	if result.Data.Task.Status == 200 {
		_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
			task.XId,
			mgz.Op().
				Set("generatedResult.url", result.Data.Task.VideoList[0].DownloadUrl).
				Set("status", "videoGenerated"),
		)
	}

	if result.Data.Task.Status == 300 {

		_, err = t.data.Mongo.Task.UpdateByIDIfExists(ctx,
			task.XId,
			mgz.Op().
				Set("status", "videoGenerateFailed"),
		)
	}

	if err != nil {
		return err
	}

	return nil
}
