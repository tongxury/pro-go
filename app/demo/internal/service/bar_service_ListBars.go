package service

import (
	"context"
	demopb "store/api/demo"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/bson"
)

func (t *BarService) ListBars(ctx context.Context, req *demopb.ListBarsRequest) (*demopb.BarList, error) {
	filter := bson.M{}
	if req.FooId != "" {
		filter["fooId"] = req.FooId
	}

	opts := &mgz.Options{
		Page: req.Page,
		Size: req.Size,
		Sort: bson.M{"createdAt": -1},
	}

	list, total, err := t.data.Mongo.Bar.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	return &demopb.BarList{
		List:  list,
		Total: total,
	}, nil
}
