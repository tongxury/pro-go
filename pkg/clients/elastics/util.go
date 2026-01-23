package elastics

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func NewMatchQuery(field, value string) *types.Query {
	return &types.Query{
		Match: map[string]types.MatchQuery{
			field: {
				Query: value,
			},
		},
	}
}

func NewMultiMatchQuery(value string, fields ...string) *types.Query {
	return &types.Query{
		MultiMatch: &types.MultiMatchQuery{
			Fields: fields,
			Query:  value,
		},
	}
}

func NewPrefixQuery(field, value string) *types.Query {
	return &types.Query{
		Prefix: map[string]types.PrefixQuery{
			field: {
				Value: value,
			},
		},
	}
}

// NewTermQuery 精准匹配查询 - 不会对查询词进行分词
func NewTermQuery(field, value string) *types.Query {
	return &types.Query{
		Term: map[string]types.TermQuery{
			field: {
				Value: value,
			},
		},
	}
}

// NewMatchPhraseQuery 精准匹配查询 - 不会对查询词进行分词
func NewMatchPhraseQuery(field, value string) *types.Query {
	return &types.Query{
		MatchPhrase: map[string]types.MatchPhraseQuery{
			field: {
				Query: value,
				Slop:  &[]int{0}[0],
			},
		},
	}
}

func NewSort(field, d string) []types.SortCombinations {
	return []types.SortCombinations{
		map[string]interface{}{
			field: map[string]string{
				"order": d,
			},
		},
	}
}

//
//// NewTermsQuery 多值精准匹配查询 - 匹配多个值中的任意一个
//func NewTermsQuery(field string, values ...interface{}) types.Query {
//	return types.Query{
//		Terms: &types.TermsQuery{
//			TermsQuery: map[string]types.TermsQueryField{
//				field: values,
//			},
//		},
//	}
//}
//
//// NewWildcardQuery 通配符查询 - 支持 * 和 ? 通配符
//func NewWildcardQuery(field, value string) types.Query {
//	return types.Query{
//		Wildcard: map[string]types.WildcardQuery{
//			field: {
//				Value: value,
//			},
//		},
//	}
//}
//
//// NewExistsQuery 字段存在性查询 - 检查字段是否存在
//func NewExistsQuery(field string) types.Query {
//	return types.Query{
//		Exists: &types.ExistsQuery{
//			Field: field,
//		},
//	}
//}
