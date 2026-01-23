package redizv2

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"store/pkg/sdk/helper"
	"time"
)

type RJSON[T Document] struct {
	keyPrefix string
	c         *RedisClient
}

func NewRJSON[T Document](c *RedisClient, keyPrefix string) *RJSON[T] {
	return &RJSON[T]{c: c, keyPrefix: keyPrefix}
}

func (t *RJSON[T]) RedisClient() *RedisClient {
	return t.c
}

func (t *RJSON[T]) Insert(ctx context.Context, values ...T) error {

	if len(values) == 0 {
		return nil
	}

	args := helper.Mapping(values, func(x T) redis.JSONSetArgs {

		return redis.JSONSetArgs{
			Key:   t.keyPrefix + x.GetKey(),
			Path:  "$",
			Value: x,
		}
	})

	_, err := t.c.JSONMSetArgs(ctx, args).Result()
	if err != nil {
		return err
	}

	return nil
}

//func (t *RJSON[T]) InsertNX(ctx context.Context, value T) error {
//	return t.c.JSONSetMode(ctx, t.keyPrefix+value.GetKey(), "$", value, "NX").Err()
//}

func (t *RJSON[T]) InsertNX(ctx context.Context, docs ...T) error {

	if len(docs) == 0 {
		return nil
	}

	pipe := t.c.Pipeline()

	for i := range docs {
		x := docs[i]

		pipe.JSONSetMode(ctx, t.keyPrefix+x.GetKey(), "$", x, "NX")

	}

	// todo 这里会插入成功也会返回 redis: nil
	_, _ = pipe.Exec(ctx)

	//log.Debugw("SetNX ing", "", "err", err, "cmds", cmds)
	//
	//if err != nil {
	//	return err
	//}

	return nil
}

type Field struct {
	Key   string
	Field string // "" 代表 root
	Val   any
	Mode  string
}

type Fields []Field

func (t *RJSON[T]) UpdateFields(ctx context.Context, fields ...Field) error {
	return t.SetFields(ctx, fields...)
}

func (t *RJSON[T]) Expire(ctx context.Context, keys []string, duration time.Duration) error {

	//}
	//
	//func (t *RedisClient) MExpire(ctx context.Context, keys []string, duration time.Duration) error {

	pipe := t.c.Pipeline()
	for i := range keys {
		pipe.Expire(ctx, t.keyPrefix+keys[i], duration)
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}

	return err
}

func (t *RJSON[T]) SetFields(ctx context.Context, fields ...Field) error {
	if len(fields) == 0 {
		return nil
	}

	pipe := t.c.Pipeline()

	for i := range fields {
		x := fields[i]

		var val any

		switch x.Val {
		case "string":
			val = String(x.Val)
		default:
			val = x.Val
		}

		field := "$." + x.Field
		if x.Field == "" {
			field = "$"
		}

		if x.Mode == "" {
			pipe.JSONSet(ctx, t.keyPrefix+x.Key, field, val)
		} else {
			pipe.JSONSetMode(ctx, t.keyPrefix+x.Key, field, val, x.Mode)
		}

	}

	// todo 这里会插入成功也会返回 redis: nil
	_, err := pipe.Exec(ctx)

	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	//if err != nil {
	//	log.Errorw("pipe.Exec err", err, "fields", fields, "cmds", helper.Mapping(cmds, func(x redis.Cmder) error {
	//		return x.Err()
	//	}))
	//}

	return nil

}

func (t *RJSON[T]) FindFieldByIds(ctx context.Context, field string, ids ...string) (map[string]any, error) {

	if len(ids) == 0 {
		return nil, nil
	}

	redisKeys := helper.Mapping(ids, func(x string) string {
		return t.keyPrefix + x
	})

	result, err := t.c.JSONMGet(ctx, field, redisKeys...).Result()
	if err != nil {
		return nil, err
	}

	mp := map[string]any{}

	for i := range ids {
		if result[i] != nil {
			mp[ids[i]] = result[i]
		}
	}

	return mp, nil
}
