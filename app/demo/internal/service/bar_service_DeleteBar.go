package service

import (
	"context"
	demopb "store/api/demo"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (t *BarService) DeleteBar(ctx context.Context, req *demopb.DeleteBarRequest) (*emptypb.Empty, error) {
	err := t.data.Mongo.Bar.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
