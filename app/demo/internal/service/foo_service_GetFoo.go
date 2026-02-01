package service

import (
	"context"
	demopb "store/api/demo"
)

func (t *FooService) GetFoo(ctx context.Context, req *demopb.GetFooRequest) (*demopb.Foo, error) {
	return t.data.Mongo.Foo.FindByID(ctx, req.Id)
}
