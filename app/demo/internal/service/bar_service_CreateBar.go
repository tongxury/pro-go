package service

import (
	"context"
	demopb "store/api/demo"
	"time"
)

func (t *BarService) CreateBar(ctx context.Context, req *demopb.CreateBarRequest) (*demopb.Bar, error) {
	bar := &demopb.Bar{
		Foo: &demopb.Foo{
			XId: req.FooId,
		},
		Status:    "active",
		CreatedAt: time.Now().Unix(),
	}

	res, err := t.data.Mongo.Bar.Insert(ctx, bar)
	if err != nil {
		return nil, err
	}

	return res, nil
}
