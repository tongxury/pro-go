package redizv2

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
)

//type RedisClusterClient struct {
//	*redis.ClusterClient
//}

type RedisClusterClient struct {
	*redis.ClusterClient
}

func NewRedisClusterClient(conf Config) *redis.ClusterClient {
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    conf.Addrs,
		Password: conf.Password,
		//ReadTimeout:  conf.ReadTimeout,
		//WriteTimeout: conf.WriteTimeout,
		Protocol: 2,
	})

	err := c.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return c
}

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(conf Config) *RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:     conf.Addrs[0],
		Password: conf.Password,
		//ReadTimeout:  conf.ReadTimeout,
		//WriteTimeout: conf.WriteTimeout,
		Protocol: 2,
	})

	err := c.Ping(context.Background()).Err()
	if err != nil {
		//panic(err)
		log.Errorw("redis client init failed err", err)
	}

	return &RedisClient{Client: c}
}
