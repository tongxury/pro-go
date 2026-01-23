package rediz

//
//import (
//	"context"
//	"fmt"
//	"github.com/redis/go-redis/v9"
//	"store/pkg/sdk/helper"
//	"store/pkg/sdk/jsonz"
//	"time"
//)
//
//type FTSearchHash[T Doc] struct {
//	*RedisClient
//	name string
//	JSON jsonz.API
//}
//
//func NewFTSearchHash[T Doc](db *RedisClient, name string) *FTSearchHash[T] {
//	return &FTSearchHash[T]{
//		RedisClient: db,
//		name:        name,
//		JSON:        jsonz.New(),
//	}
//}
//
//type Query struct {
//	Field   string
//	Value interface{}
//	Type  string // tag numeric text
//}
//
//func (q Query) String() string {
//	switch q.Type {
//	case "numeric":
//	case "text":
//		return fmt.Sprintf("@%s: %v", q.Field, q.Value)
//	default:
//		return fmt.Sprintf("@%s: %v", q.Field, q.Value)
//	}
//
//	return ""
//}
//
//type AndQueries []Query
//
//type FindArgs struct {
//	AndQueries AndQueries
//}
//
//// Find https://redis.io/docs/latest/develop/interact/search-and-query/query/exact-match/
//func (t *FTSearchHash[T]) Find(ctx context.Context, args FindArgs) ([]*T, error) {
//
//	return nil, nil
//}
//
//func (t *FTSearchHash[T]) FindFieldByIds(ctx context.Context, field string, ids ...string) (map[string]any, error) {
//
//	cmds, err := t.RedisClient.Pipelined(ctx, func(pipe redis.Pipeliner) error {
//
//		for _, id := range ids {
//			pipe.HGet(ctx, t.RedisKey(id), field)
//		}
//
//		return nil
//	})
//
//	if err != nil {
//		return nil, err
//	}
//
//	fmt.Println(cmds)
//
//	return nil, nil
//}
//
//func (t *FTSearchHash[T]) FindAllById(ctx context.Context, id string) (*T, error) {
//
//	result, err := t.HGetAll(ctx, t.RedisKey(id)).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	var val T
//	err = t.JSON.MapToStruct(result, &val)
//	if err != nil {
//		return nil, err
//	}
//
//	return &val, nil
//}
//
//type Field struct {
//	Field string
//	Val any
//}
//
//type Fields []Field
//
//func (t *FTSearchHash[T]) SetFields(ctx context.Context, args map[string]Fields) error {
//
//	if len(args) == 0 {
//		return nil
//	}
//
//	return nil
//}
//
//func (t *FTSearchHash[T]) Set(ctx context.Context, docs ...T) error {
//
//	if len(docs) == 0 {
//		return nil
//	}
//
//	pipe := t.Pipeline()
//
//	for i := range docs {
//		x := docs[i]
//
//		pipe.HSet(ctx, t.RedisKey(x.Field()), x)
//	}
//
//	_, err := pipe.Exec(ctx)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}
//
//type SetArgs[T Doc] struct {
//	Doc         T
//	OverrideAll bool
//}
//
//func (t *FTSearchHash[T]) RedisKey(id string) string {
//	return fmt.Sprintf("hash.%s:%s", t.name, id)
//	//return fmt.Sprintf("stashedDexTrade:%s", id)
//}
//
//func (t *FTSearchHash[T]) indexKey() string {
//	return fmt.Sprintf("idx:%s", t.name)
//}
//
//func (t *FTSearchHash[T]) Expires(ctx context.Context, ids []string, duration time.Duration) error {
//
//	if len(ids) == 0 {
//		return nil
//	}
//
//	keys := helper.Mapping(ids, func(x string) string {
//		return t.RedisKey(x)
//	})
//
//	return t.MExpire(ctx, keys, duration)
//}
