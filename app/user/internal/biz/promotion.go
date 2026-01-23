package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	typepb "store/api/user/types"
	"store/app/user/internal/data"
)

type PromotionBiz struct {
	data *data.Data
}

func NewPromotionBiz(data *data.Data) *PromotionBiz {
	return &PromotionBiz{data: data}
}

func (t *PromotionBiz) IsAvailable(ctx context.Context, code, level, cycle string) (bool, error) {

	return false, nil
}

func (t *PromotionBiz) GetById(ctx context.Context, id string) (*typepb.Promotion, error) {

	promotion, found := t.data.BizConfig.Promotions.AsMap()[id]
	if !found {
		return nil, errors.BadRequest("invalid code: "+id, "")
	}

	used, err := t.data.Repos.Promotion.GetUsed(ctx, id)
	if err != nil {
		return nil, err
	}

	return &typepb.Promotion{
		Id:    promotion.Id,
		Limit: promotion.Limit,
		Used:  used,
	}, nil
}

func (t *PromotionBiz) IsValid(ctx context.Context, id string) (bool, error) {

	prom, err := t.GetById(ctx, id)
	if err != nil {
		return false, err
	}

	if prom == nil {
		return false, nil
	}

	return prom.Used < prom.Limit, nil
}
