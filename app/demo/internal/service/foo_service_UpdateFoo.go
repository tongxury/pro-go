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
	foo.Name = req.Name
	foo.Avatar = req.Avatar
	foo.Desc = req.Desc
	foo.UpdatedAt = time.Now().Unix()

	_, err = t.data.Mongo.Foo.ReplaceByID(ctx, req.Id, foo)
	if err != nil {
		return nil, err
	}

	return foo, nil
}
