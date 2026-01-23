package service

import (
	"context"
	projpb "store/api/proj"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (t ProjService) UpdateCommodity(ctx context.Context, request *projpb.UpdateCommodityRequest) (*emptypb.Empty, error) {

	switch request.Action {
	case "delete":
		err := t.data.Mongo.Commodity.DeleteByID(ctx, request.Id)
		if err != nil {
			return nil, err
		}
	default:

	}

	return &emptypb.Empty{}, nil
}
