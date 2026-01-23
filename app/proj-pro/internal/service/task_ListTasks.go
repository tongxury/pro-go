package service

import (
	"context"
	"encoding/json"
	projpb "store/api/proj"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"go.mongodb.org/mongo-driver/bson"
)

func (t ProjService) ListTasks(ctx context.Context, params *projpb.ListTaskRequest) (*projpb.TaskList, error) {

	userId := krathelper.RequireUserId(ctx)

	filter := bson.M{"userId": userId}

	if params.Keyword != "" {
		filter["$or"] = []bson.M{
			{"commodity.brand": bson.M{"$regex": params.Keyword, "$options": "i"}},
			{"commodity.title": bson.M{"$regex": params.Keyword, "$options": "i"}},
			{"commodity.category": bson.M{"$regex": params.Keyword, "$options": "i"}},
		}
	}

	list, err := t.data.Mongo.Task.List(ctx,
		filter,
		mgz.Find().
			SetFields(params.ReturnFields).
			PageSize(params.Page, params.Size).
			SetSort("createdAt", -1).
			B(),

		//mongoz.Paging(params.Page, 10).
		//	SetSort(bson.M{"createdAt": -1}).
		//	SetProjection(bson.M{
		//		"commodity":    1,
		//		"targetChance": 1,
		//		"_id":          1,
		//	}),
	)
	if err != nil {
		return nil, err
	}

	return &projpb.TaskList{
		List: list,
	}, err
}

func (t ProjService) parseToItems(hits []types.Hit) []*projpb.Item {
	var items []*projpb.Item
	for _, x := range hits {

		var y *projpb.Item
		_ = json.Unmarshal(x.Source_, &y)

		y.Id = *x.Id_

		items = append(items, y)
	}

	return items
}
