package service

import (
	"context"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/third/bytedance/vikingdb"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjAdminService) SyncTemplateSegmentsToVikingDB() {

	ctx := context.Background()

	list, err := t.data.Mongo.TemplateSegment.List(ctx,
		//mongoz.NotExists("extra.synced"),
		bson.M{
			"status":       bson.M{"$regex": "^completed"},
			"extra.synced": bson.M{"$exists": false},
		},
	)
	if err != nil {
		log.Errorw("List err", err)
		return
	}

	if len(list) == 0 {
		return
	}

	for _, x := range list {

		cname := x.Root.GetCommodity().GetName()

		if cname != "" {
			_, err = t.data.VikingDB.Upsert(ctx, vikingdb.UpsertRequest{
				Collection: "segment_commodity_coll",
				Data: map[string]any{
					"id":          x.XId,
					"description": x.Root.GetCommodity().GetName(),
				},
			})
			if err != nil {
				log.Errorw("SyncTemplateSegmentsToVikingDB err", err, "id", x.XId, "description", x.Root.GetCommodity().GetName())
				continue
			}

		}

		tags := x.GetTags()
		if len(tags) > 0 {
			_, err = t.data.VikingDB.Upsert(ctx, vikingdb.UpsertRequest{
				Collection: "segment_video_coll",
				Data: map[string]any{
					"id":          x.XId,
					"description": strings.Join(tags, ","),
				},
			})
			if err != nil {
				log.Errorw("SyncTemplateSegmentsToVikingDB err", err, "id", x.XId, "description", strings.Join(x.GetTags(), ","))
				continue
			}
		}
		
		_, err = t.data.Mongo.TemplateSegment.UpdateByIDXX(ctx, x.XId, mgz.Set(bson.M{"extra.synced": true}))
		if err != nil {
			return
		}

	}
}
