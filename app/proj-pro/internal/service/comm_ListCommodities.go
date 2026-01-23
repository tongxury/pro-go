package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"store/pkg/sdk/third/bytedance/tos"

	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) ListCommodities(ctx context.Context, params *projpb.ListCommodityRequest) (*projpb.CommodityList, error) {

	userId := krathelper.RequireUserId(ctx)

	filter := bson.M{"userId": userId}

	if params.Keyword != "" {
		filter["$or"] = []bson.M{
			{"brand": bson.M{"$regex": params.Keyword, "$options": "i"}},
			{"title": bson.M{"$regex": params.Keyword, "$options": "i"}},
			{"category": bson.M{"$regex": params.Keyword, "$options": "i"}},
		}
	}

	//

	if len(params.Ids) > 0 {
		filter["_id"] = bson.M{"$in": mgz.ObjectIds(params.Ids)}
	}

	list, err := t.data.Mongo.Commodity.List(ctx, filter,
		mgz.Find().
			Paging(params.Page, params.Size).
			SetFields(params.ReturnFields).
			SetSort("createdAt", -1).
			B(),
	)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		//	准备示例数据
		demoList, err := t.data.Mongo.Commodity.List(ctx, bson.M{},
			mgz.Find().
				Paging(1, 2).
				SetSort("createdAt", -1).
				B(),
		)
		if err != nil {
			return nil, err
		}

		for _, x := range demoList {
			x.UserId = userId
		}

		_, err = t.data.Mongo.Commodity.InsertMany(ctx, demoList...)
		if err != nil {
			return nil, err
		}

		list = demoList
	}

	for i := range list {
		list[i].Images = tos.ChangeMany(list[i].Images[:1])
	}

	return &projpb.CommodityList{
		List: list,
	}, err
}
