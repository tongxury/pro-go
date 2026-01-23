package rediz

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"time"
)

type FTSearchJSON[T Doc] struct {
	*RedisClient
	name string
	//JSON jsonz.API
}

func NewFTSearchJSON[T Doc](db *RedisClient, name string) *FTSearchJSON[T] {
	return &FTSearchJSON[T]{
		RedisClient: db,
		name:        name,
		//JSON:        jsonz.New(),
	}
}

//type OrQueries []Query
//
//func (qs OrQueries) String() string {
//
//	if len(qs) == 0 {
//		return "*"
//	}
//
//	var queries []string
//	for i := range qs {
//
//		x := qs[i]
//
//		queries = append(queries, fmt.Sprintf("(@%s: %v)", x.Field, x.ValueExpr))
//	}
//
//	return strings.Join(queries, " | ")
//}

type FindArgs struct {
	Index   string
	Queries Q
	SortBy  []redis.FTSearchSortBy
	Return  []redis.FTSearchReturn
	Limit   int
	Offset  int
}

func (t *FTSearchJSON[T]) FindFields(ctx context.Context, args FindArgs) ([]map[string]string, error) {

	query := args.Queries.String()

	index := t.indexKey()

	ftResult, err := t.FTSearchWithArgs(ctx,
		index,
		query,
		&redis.FTSearchOptions{
			Return:      args.Return,
			SortBy:      args.SortBy,
			LimitOffset: args.Offset,
			Limit:       args.Limit,
		}).Result()
	if err != nil {
		return nil, err
	}

	var results []map[string]string
	for i := range ftResult.Docs {
		x := ftResult.Docs[i]

		results = append(results, x.Fields)
	}

	return results, nil
}

// Find https://redis.io/docs/latest/develop/interact/search-and-query/query/exact-match/
func (t *FTSearchJSON[T]) Find(ctx context.Context, args FindArgs) ([]T, error) {

	query := args.Queries.String()

	ftResult, err := t.FTSearchWithArgs(ctx,
		helper.OrString(args.Index, t.indexKey()),
		query,
		&redis.FTSearchOptions{
			//NoContent:       false,
			//Verbatim:        false,
			//NoStopWords:     false,
			//WithScores:      false,
			//WithPayloads:    false,
			//WithSortKeys:    false,
			//Filters:         nil,
			//GeoFilter:       nil,
			//InKeys:          nil,
			//InFields:        nil,
			Return: args.Return,
			//Slop:            0,
			//Timeout:         0,
			//InOrder:         false,
			//Language:        "",
			//Expander:        "",
			//Scorer:          "",
			//ExplainScore:    false,
			//Payload:         "",
			SortBy: args.SortBy,
			//SortByWithCount: false,
			LimitOffset: args.Offset,
			Limit:       args.Limit,
			//Params:         nil,
			//DialectVersion: ,
		}).Result()
	if err != nil {
		return nil, err
	}

	var results []T
	for i := range ftResult.Docs {
		x := ftResult.Docs[i]

		if val, found := x.Fields["$"]; found {

			var tmp T
			_ = conv.J2S([]byte(val), &tmp)

			results = append(results, tmp)
		}
	}

	return results, nil
}

func (t *FTSearchJSON[T]) FindFieldByIds(ctx context.Context, field string, ids ...string) (map[string]any, error) {

	if len(ids) == 0 {
		return nil, nil
	}

	redisKeys := helper.Mapping(ids, func(x string) string {
		return t.RedisKey(x)
	})

	result, err := t.JSONMGet(ctx, field, redisKeys...).Result()
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

func (t *FTSearchJSON[T]) FindAllById(ctx context.Context, id string) (*T, error) {

	result, err := t.JSONGet(ctx, t.RedisKey(id)).Result()
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, nil
	}

	var rp T
	err = conv.J2S([]byte(result), &rp)
	if err != nil {
		return nil, err
	}

	return &rp, nil
}

func (t *FTSearchJSON[T]) FindFieldsById(ctx context.Context, id string, paths ...string) (*T, error) {

	var pathList []string
	for _, x := range paths {
		pathList = append(pathList, "$."+x)
	}

	result, err := t.JSONGet(ctx, t.RedisKey(id), pathList...).Result()
	if err != nil {
		return nil, err
	}

	if result == "" {
		return nil, nil
	}

	var rp T
	err = conv.J2S([]byte(result), &rp)
	if err != nil {
		return nil, err
	}

	return &rp, nil
}

type Field struct {
	Key   string
	Field string
	Val   any
	Mode  string
}

type Fields []Field

func (t *FTSearchJSON[T]) UpdateFields(ctx context.Context, fields ...Field) error {
	return t.SetFields(ctx, fields...)
}

func (t *FTSearchJSON[T]) SetFields(ctx context.Context, fields ...Field) error {
	if len(fields) == 0 {
		return nil
	}

	pipe := t.Pipeline()

	for i := range fields {
		x := fields[i]

		switch x.Val {
		case "string":
			pipe.JSONSetMode(ctx, t.RedisKey(x.Key), "$."+x.Field, String(x.Val), x.Mode)
		default:
			pipe.JSONSetMode(ctx, t.RedisKey(x.Key), "$."+x.Field, x.Val, x.Mode)
		}

	}

	// todo 这里会插入成功也会返回 redis: nil
	cmds, err := pipe.Exec(ctx)

	if err != nil {
		log.Debugw("SetFields ing", "", "err", err, "cmds", helper.Filter(cmds, func(param redis.Cmder) bool {
			return !errors.Is(param.Err(), redis.Nil)
		}))
	}

	return nil

}

func (t *FTSearchJSON[T]) Replace(ctx context.Context, docs ...T) error {
	if len(docs) == 0 {
		return nil
	}

	var params []interface{}

	for _, x := range docs {
		redisKey := t.RedisKey(x.Key())

		marshal := conv.S2J(x)

		params = append(params, redisKey, "$", marshal)
	}

	return t.JSONMSet(ctx, params...).Err()
}

func (t *FTSearchJSON[T]) SetNX(ctx context.Context, docs ...T) error {

	if len(docs) == 0 {
		return nil
	}

	pipe := t.Pipeline()

	for i := range docs {
		x := docs[i]

		pipe.JSONSetMode(ctx, t.RedisKey(x.Key()), "$", x, "NX")

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

type SetArgs[T Doc] struct {
	Doc         T
	OverrideAll bool
}

func (t *FTSearchJSON[T]) RedisKey(id string) string {
	return fmt.Sprintf("rj.%s:%s", t.name, id)
	//return fmt.Sprintf("stashedDexTrade:%s", id)
}

func (t *FTSearchJSON[T]) indexKey() string {
	return fmt.Sprintf("idx:%s", t.name)
}

func (t *FTSearchJSON[T]) Expires(ctx context.Context, ids []string, duration time.Duration) error {

	if len(ids) == 0 {
		return nil
	}

	keys := helper.Mapping(ids, func(x string) string {
		return t.RedisKey(x)
	})

	if err := t.MExpire(ctx, keys, duration); err != nil {
		return t.MExpire(ctx, keys, duration)
	}
	return nil
}
