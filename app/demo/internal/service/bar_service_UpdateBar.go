package service

import (
	"context"
	demopb "store/api/demo"
)

func (t *BarService) UpdateBar(ctx context.Context, req *demopb.UpdateBarRequest) (*demopb.Bar, error) {
	bar, err := t.data.Mongo.Bar.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if bar == nil {
		return nil, nil
	}

	bar.Status = req.Status
	// AI 提示：映射扩展配置
	// bar.Config = req.Config

	_, err = t.data.Mongo.Bar.ReplaceByID(ctx, req.Id, bar)
	if err != nil {
		return nil, err
	}

	return bar, nil
}
