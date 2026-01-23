package rediz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

//type RedisClusterClient struct {
//	*redis.ClusterClient
//}

type RedisClient struct {
	*redis.Client
}

func (t *RedisClient) MExpire(ctx context.Context, keys []string, duration time.Duration) error {

	pipe := t.Pipeline()
	for i := range keys {
		pipe.Expire(ctx, keys[i], duration)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return err
}

func NewRedisClient(conf Config) *RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     conf.Addrs[0],
		Password: conf.Password,
		//ReadTimeout:  conf.ReadTimeout,
		//WriteTimeout: conf.WriteTimeout,
		DB:       conf.DB,
		Protocol: 2,
	})
	//
	err := c.Ping(context.Background()).Err()
	if err != nil {
		log.Debugw("redis client init failed", conf, "err", err)
		panic(err)
	}

	return &RedisClient{Client: c}
}

//func NewRedisClusterClient(conf confcenter.Redis) *RedisClusterClient {
//	c := redis.NewClusterClient(&redis.ClusterOptions{
//		Addrs:        conf.Addrs,
//		Password:     conf.Password,
//		ReadTimeout:  conf.ReadTimeout,
//		WriteTimeout: conf.WriteTimeout,
//	})
//
//	err := c.Ping(context.Background()).Err()
//	if err != nil {
//		panic(err)
//	}
//
//	return &RedisClusterClient{ClusterClient: c}
//}
