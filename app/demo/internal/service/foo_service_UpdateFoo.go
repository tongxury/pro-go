package service

import (
	"context"
	demopb "store/api/demo"
	"time"
)

func (t *FooService) UpdateFoo(ctx context.Context, req *demopb.UpdateFooRequest) (*demopb.Foo, error) {
	foo, err := t.data.Mongo.Foo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if foo == nil {
		return nil, nil // Or handle as error
	}

	// AI 提示：根据需求更新字段
	foo.Foo = req.Name
	foo.UpdatedAt = time.Now().Unix()

	err = t.data.Mongo.Foo.UpdateByID(ctx, req.Id, foo)
	if err != nil {
		return nil, err
	}

	return foo, nil
}
