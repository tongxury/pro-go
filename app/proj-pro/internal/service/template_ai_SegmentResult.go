package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper/wg"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) SegmentResult() {
	//t.data.Elastics.Create()
	ctx := context.Background()

	hits, err := t.data.Mongo.Template.List(ctx, bson.M{"status": "segmented"})
	if err != nil {
		log.Errorw("Search err", err, "index", "items")
		return
	}

	if len(hits) == 0 {
		return
	}

	wg.WaitGroup(ctx, hits, t.segmentResult)
}

func (t ProjService) segmentResult(ctx context.Context, tpl *projpb.Resource) error {
	//logger := log.NewHelper(log.With(log.DefaultLogger,
	//	"func", "segmentResult",
	//	"item", tpl.XId,
	//))
	//
	//logger.Debugw("start segmentResult", "")

	list, err := t.data.Mongo.TemplateSegment.List(ctx, bson.M{"root._id": tpl.XId})
	if err != nil {
		return err
	}

	if len(list) == 0 {
		t.data.Mongo.Template.UpdateByIDIfExists(ctx, tpl.XId, mgz.Op().Set("status", "created"))
		return nil
	}

	for _, x := range list {
		if x.Status != "completed" {
			return nil
		}
	}

	t.data.Mongo.Template.UpdateByIDIfExists(ctx, tpl.XId, mgz.Op().Set("status", "completed"))

	return nil
}
