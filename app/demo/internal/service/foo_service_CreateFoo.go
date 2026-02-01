package service

import (
	"context"
	demopb "store/api/demo"
	"time"
)

func (t *FooService) CreateFoo(ctx context.Context, req *demopb.CreateFooRequest) (*demopb.Foo, error) {
	foo := &demopb.Foo{
		Foo:       req.Name, // 假设映射
		Status:    "active",
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	res, err := t.data.Mongo.Foo.Insert(ctx, foo)
	if err != nil {
		return nil, err
	}

	return res, nil
}
