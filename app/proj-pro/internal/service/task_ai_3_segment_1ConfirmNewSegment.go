package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	projpb "store/api/proj"
	"store/pkg/sdk/helper/wg"
	"strings"
)

func (t ProjService) ConfirmNewSegment() {

	ctx := context.Background()

	list, err := t.data.Mongo.TaskSegment.List(ctx, bson.M{"status": "textGenerated"})
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	wg.WaitGroup(ctx, list, t.confirmNewSegment)
}

func (t ProjService) confirmNewSegment(ctx context.Context, taskSegment *projpb.TaskSegment) error {

	log.Debugw("doGenerateV2", "confirmNewSegment", "taskSegment", taskSegment.XId)

	segments, err := t.data.GrpcClients.ProjAdminClient.SearchMatchedItemSegments(ctx, &projpb.SearchMatchedItemSegmentsParams{
		Keyword:          strings.Join(taskSegment.Segment.Tags, ","),
		CommodityKeyword: taskSegment.Task.Commodity.GetTitle(),
		SubtitleKeyword:  taskSegment.Subtitle,
	})
	if err != nil {
		return err
	}

	if len(segments.List) > 0 {
		_, err := t.data.Mongo.TaskSegment.UpdateOneXXById(ctx,
			taskSegment.XId,
			bson.M{
				"newSegment": segments.List[0],
				"status":     "newSegmentConfirmed",
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
