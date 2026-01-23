package service

import (
	"context"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/third/bytedance/vikingdb"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjAdminService) SyncTemplatesToVikingDB() {

	ctx := context.Background()

	list, err := t.data.Mongo.Template.List(ctx,
		//mongoz.And(
		//	mongoz.NotExists("extra.synced"),
		//),
		//
		//mongoz.NotExists("extra.synced"),
		bson.M{
			//"userId":       bson.M{"$exists": false},
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

	x := list[0]

	if x.UserId == "" {
		x.UserId = "system"
	}

	_, err = t.data.VikingDB.Upsert(ctx, vikingdb.UpsertRequest{
		Collection: "template_commodity_coll",
		Data: map[string]any{
			"id":          x.XId,
			"userId":      x.UserId,
			"description": x.GetCommodity().GetName(),
		},
	})
	if err != nil {
		log.Errorw("SyncTemplateSegmentsToVikingDB err", err)
		return
	}

	_, err = t.data.Mongo.Template.UpdateByIDXX(ctx, x.XId, mgz.Set(bson.M{"extra.synced": true}))
	if err != nil {
		return
	}
}
