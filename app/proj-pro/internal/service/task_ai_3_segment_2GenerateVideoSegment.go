package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/sdk/helper/wg"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) GenerateVideoSegment() {

	ctx := context.Background()

	list, err := t.data.Mongo.TaskSegment.List(ctx, bson.M{"status": "aiVideoSegmentGenerating"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.generateVideoSegment3)
}

func (t ProjService) generateVideoSegment3(ctx context.Context, taskSegment *projpb.TaskSegment) error {

	//	assetId := taskSegment.AssetId
	//	if taskSegment.AssetId != "" {
	//
	//		asset, err := t.data.Mongo.Asset.FindByID(ctx, assetId)
	//		if err != nil {
	//			return err
	//		}
	//
	//		if asset.GetStatus() == "completed" {
	//
	//			generatedTasks := []*projpb.GeneratedTask{
	//				{
	//					//Id:       video.ID,
	//					VideoUrl: asset.GetUrl(),
	//				},
	//			}
	//
	//			_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, taskSegment.XId, bson.M{
	//				"generatedTasks": generatedTasks,
	//				"status":         "aiVideoSegmentGenerated",
	//			})
	//			if err != nil {
	//				return err
	//			}
	//		}
	//
	//		if asset.GetStatus() == "failed" {
	//			_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, taskSegment.XId, bson.M{
	//				"status": "aiVideoSegmentGenerateFailed",
	//			})
	//			if err != nil {
	//				return err
	//			}
	//		}
	//
	//		return nil
	//	}
	//
	//	// 检测mode
	//	task, err := t.data.Mongo.Task.GetById(ctx, taskSegment.TaskId)
	//	if err != nil {
	//		return err
	//	}
	//
	//	if task.Mode == "sequential" {
	//
	//	}
	//
	//	//
	//	newAsset, err := t.data.Mongo.Asset.Insert(ctx, &projpb.Asset{
	//		Commodity: taskSegment.Task.Commodity,
	//		Segment:   taskSegment.Segment,
	//		UserId:    taskSegment.Task.UserId,
	//		Status:    "created",
	//		CreatedAt: time.Now().Unix(),
	//		Prompt: fmt.Sprintf(`
	//- **视频的口播文案要用这个: %s**
	//
	//- %s
	//`, taskSegment.Subtitle,
	//			taskSegment.Prompts["generateVideoSegment"]),
	//	})
	//	if err != nil {
	//		return err
	//	}
	//
	//	_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, taskSegment.XId, bson.M{
	//		"extra.context.assetId": newAsset.XId,
	//		"status":                "aiVideoSegmentGenerating",
	//	})
	//
	//	if err != nil {
	//		return err
	//	}

	return nil
}
