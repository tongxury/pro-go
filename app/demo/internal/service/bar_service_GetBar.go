package service

import (
	"context"
	demopb "store/api/demo"
)

func (t *BarService) GetBar(ctx context.Context, req *demopb.GetBarRequest) (*demopb.Bar, error) {
	return t.data.Mongo.Bar.FindByID(ctx, req.Id)
}
