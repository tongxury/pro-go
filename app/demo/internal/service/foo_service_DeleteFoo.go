package service

import (
	"context"
	demopb "store/api/demo"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *FooService) DeleteFoo(ctx context.Context, req *demopb.DeleteFooRequest) (*emptypb.Empty, error) {
	err := t.data.Mongo.Foo.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
