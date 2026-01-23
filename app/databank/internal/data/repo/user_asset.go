package repo

import (
	"context"
	"store/app/databank/internal/data/repo/ent"
	"store/pkg/clients"
)

type UserAssetRepo struct {
	db    *ent.Client
	redis *clients.RedisClient
}

func NewUserAssetRepo(db *ent.Client, redis *clients.RedisClient) *UserAssetRepo {
	return &UserAssetRepo{db: db, redis: redis}
}

func (t *UserAssetRepo) Insert(ctx context.Context, userId, sessionId string) error {
	return nil
}
