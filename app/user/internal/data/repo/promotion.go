package repo

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"store/app/user/internal/data/repo/ent"
	"store/app/user/internal/data/repo/ent/promotion"
	"store/pkg/clients"
)

type PromotionRepo struct {
	db    *ent.Client
	redis *clients.RedisClient
}

func NewPromotionRepo(db *ent.Client, redis *clients.RedisClient) *PromotionRepo {
	return &PromotionRepo{db: db, redis: redis}
}

func (t *PromotionRepo) FindByCode(ctx context.Context, code string) ([]*ent.Promotion, error) {

	all, err := t.db.Promotion.Query().Where(promotion.Code(code)).All(ctx)
	if err != nil {
		return nil, err
	}

	return all, nil
}

func (t *PromotionRepo) GetUsed(ctx context.Context, id string) (int64, error) {

	key := fmt.Sprintf("promotion.used:%s", id)
	limit, err := t.redis.Get(ctx, key).Int64()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}

	return limit, nil
}
