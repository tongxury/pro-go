package data

import (
	"context"
	"errors"
	"store/pkg/rediz"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	rc *rediz.RedisClient
}

func (t *RedisCache) Get(ctx context.Context, key string) (string, error) {
	result, err := t.rc.Get(ctx, key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return "", err
	}

	return result, nil
}

func (t *RedisCache) Set(ctx context.Context, key string, value string, dur time.Duration) error {
	return t.rc.Set(ctx, key, value, dur).Err()
}

func NewRedisCache(rc *rediz.RedisClient) *RedisCache {
	return &RedisCache{
		rc: rc,
	}
}
