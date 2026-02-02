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

	opts := mgz.Find().
		Paging(req.Page, req.Size).
		SetSort("createdAt", -1).
		B()

	list, total, err := t.data.Mongo.Bar.ListAndCount(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	return &demopb.BarList{
		List:  list,
		Total: total,
	}, nil
}
