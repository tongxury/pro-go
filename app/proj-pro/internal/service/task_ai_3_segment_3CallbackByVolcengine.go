package service

import (
	"context"
	"io"
	projpb "store/api/proj"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/golang/protobuf/ptypes/empty"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) CallbackByVolcengine(ctx context.Context, params *projpb.CallbackByVolcengineParams) (*empty.Empty, error) {

	req, _ := http.RequestFromServerContext(ctx)

	all, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	log.Debug("CallbackByVolcengine", string(all))

	if params.Status != "succeeded" {
		return &empty.Empty{}, nil
	}

	_, err = t.data.Mongo.TaskSegment.UpdateOneXX(ctx,
		bson.M{"generatedTasks.id": params.Id},
		bson.M{
			"generatedTasks.$.videoUrl":     params.Content.VideoUrl,
			"generatedTasks.$.lastFrameUrl": params.Content.LastFrameUrl,
			"status":                        "aiVideoSegmentGenerated",
		})
	if err != nil {
		log.Error("CallbackByVolcengine", err)
		return nil, err
	}

	return &empty.Empty{}, nil
}
