package clients

import (
	"context"
	"github.com/redis/go-redis/v9"
	"store/pkg/clients/redizv2"
	"store/pkg/confcenter"
)

type RedisClusterClient struct {
	*redis.ClusterClient
}

type RedisClient struct {
	*redis.Client
}

func NewRedisClient(conf confcenter.Redis) *RedisClient {
	c := redis.NewClient(&redis.Options{
		Addr:         conf.Addrs[0],
		Password:     conf.Password,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		Protocol:     2,
	})

	err := c.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return &RedisClient{Client: c}
}

func NewRedisClusterClient(conf redizv2.ClusterConfig) *RedisClusterClient {
	c := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    conf.Addrs,
		Password: conf.Password,
		//ReadTimeout:  conf.ReadTimeout,
		//WriteTimeout: conf.WriteTimeout,
	})

	err := c.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return &RedisClusterClient{ClusterClient: c}
}
