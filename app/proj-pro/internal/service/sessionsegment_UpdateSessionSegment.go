package service

import (
	"context"
	"fmt"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *SessionService) UpdateSessionSegment(ctx context.Context, params *projpb.UpdateSessionSegmentRequest) (*empty.Empty, error) {

	fmt.Println("UpdateSessionSegment")

	var err error
	switch params.Action {
	case "deleteAsset":
		_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id,
			bson.M{
				"status":  "newSegmentConfirmed",
				"assetId": "",
			})

	case "updateAsset":

		assetId := params.Params["assetId"]

		asset, err := t.data.Mongo.Asset.GetById(ctx, assetId)
		if err != nil {
			return nil, err
		}

		_, err = t.data.Mongo.SessionSegment.UpdateByIDIfExists(ctx,
			params.Id,
			mgz.Op().Set("asset", asset),
		)

	case "generateVideoSegment":
		_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id,
			bson.M{
				"status":                "aiVideoSegmentGenerating",
				"extra.context.assetId": "",
			})

	case "updateStatus":
		_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id, bson.M{"status": params.Params["status"]})

	case "update":
		_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id, bson.M{params.Params["field"]: params.Params["value"]})

	case "generateAISegment":
		_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id, bson.M{
			"status":             "aiSegmentGenerating",
			"frames.firstOrigin": params.Params["first"],
			"frames.lastOrigin":  params.Params["last"],
		})

	case "updateAISegment":

		cat := params.Params["cat"]

		if cat == "last" {
			_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id,
				bson.M{"frames.last.$[].selected": false},
			)

			_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id,
				bson.M{fmt.Sprintf("frames.last.%s.selected", params.Params["index"]): true})

		} else {
			_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id,
				bson.M{"frames.first.$[].selected": false},
			)

			_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id,
				bson.M{fmt.Sprintf("frames.first.%s.selected", params.Params["index"]): true})
		}

	case "addAISegment":

		_, err = t.data.Mongo.TaskSegment.UpdateOneXXPushById(ctx, params.Id, "frames.first",
			[]bson.M{{"url": params.Params["url"]}}, 0)

	case "updateAIVideoSegment":
		_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id,
			bson.M{"generatedTasks.$[].selected": false},
		)
		_, err = t.data.Mongo.TaskSegment.UpdateOneXXById(ctx, params.Id,
			bson.M{fmt.Sprintf("generatedTasks.%s.selected", params.Params["index"]): true})
	}

	if err != nil {
		log.Errorw("update task segment failed error", err, "params", params)
		return nil, err
	}

	return &empty.Empty{}, nil
}
