package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) GetTask(ctx context.Context, params *projpb.GetTaskRequest) (*projpb.Task, error) {

	task, err := t.data.Mongo.Task.GetById(ctx, params.Id,
		mgz.Find().SetFields(params.ReturnFields).B())
	if err != nil {
		return nil, err
	}

	segments, err := t.data.Mongo.TaskSegment.List(ctx, bson.M{"taskId": task.XId},
		mgz.Find().SetFields(params.ReturnFields).B())
	if err != nil {
		return nil, err
	}

	// 补充 asset信息
	var assetIds []string
	for _, x := range segments {
		if x.AssetId != "" {
			assetIds = append(assetIds, x.AssetId)
		}
	}

	if len(assetIds) > 0 {
		assets, err := t.data.Mongo.Asset.List(ctx,
			mgz.Filter().Ids(assetIds).B(),
			mgz.Find().SetFields(params.ReturnFields).B(),
		)
		if err != nil {
			return nil, err
		}

		assetsMap := map[string]*projpb.Asset{}
		for _, x := range assets {
			assetsMap[x.XId] = x
		}

		for _, x := range segments {
			x.Asset = assetsMap[x.AssetId]
		}
	}

	task.Segments = segments

	return task, nil
}
