package redizv2

import (
	"context"
	"github.com/redis/go-redis/v9"
	"store/pkg/sdk/conv"
)

type FTSearch[T Document] struct {
	c *RedisClient
}

func NewFTSearch[T Document](c *RedisClient) *FTSearch[T] {
	return &FTSearch[T]{
		c: c,
	}
}

type FindArgs struct {
	Index   string
	Queries Q
	SortBy  []redis.FTSearchSortBy
	Return  []redis.FTSearchReturn
	Limit   int
	Offset  int
}

//func (t *FTSearch[T]) FindFields(ctx context.Context, args FindArgs) ([]map[string]string, error) {
//
//	query := args.Queries.String()
//
//	ftResult, err := t.c.FTSearchWithArgs(ctx,
//		args.Index,
//		query,
//		&redis.FTSearchOptions{
//			Return:      args.Return,
//			SortBy:      args.SortBy,
//			LimitOffset: args.Offset,
//			Limit:       args.Limit,
//		}).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	var results []map[string]string
//	for i := range ftResult.Docs {
//		x := ftResult.Docs[i]
//
//		results = append(results, x.Fields)
//	}
//
//	return results, nil
//}

// Find https://redis.io/docs/latest/develop/interact/search-and-query/query/exact-match/
func (t *FTSearch[T]) Find(ctx context.Context, args FindArgs) ([]T, int64, error) {

	query := args.Queries.String()

	ftResult, err := t.c.FTSearchWithArgs(ctx,
		args.Index,
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
		return nil, 0, err
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

	return results, int64(ftResult.Total), nil
}
