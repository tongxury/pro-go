package service

import (
	"context"
	projpb "store/api/proj"
)

func (t ProjService) GetCommodity(ctx context.Context, params *projpb.GetCommodityRequest) (*projpb.Commodity, error) {

	task, err := t.data.Mongo.Commodity.FindByID(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	return task, nil
}
