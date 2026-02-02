package service

import (
	"context"
	demopb "store/api/demo"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/bson"
)

func (t *FooService) ListFoos(ctx context.Context, req *demopb.ListFoosRequest) (*demopb.FooList, error) {
	filter := bson.M{}
	if req.Keyword != "" {
		filter["foo"] = bson.M{"$regex": req.Keyword, "$options": "i"}
	}

	opts := mgz.Find().
		Paging(req.Page, req.Size).
		SetSort("createdAt", -1).
		B()

	list, total, err := t.data.Mongo.Foo.ListAndCount(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	return &demopb.FooList{
		List:  list,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}
