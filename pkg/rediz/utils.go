package rediz

import (
	"fmt"
	"strings"
)

var empty = "__EMPTY__"

//func IsTagEmpty(field string) Query {
//	return Query{
//		Field:       field,
//		ValueExpr: "{" + empty + "}",
//	}
//}
//
//func Empty(src string) string {
//	if src == "" {
//		return empty
//	}
//
//	return src
//}

func String(val any) string {
	return fmt.Sprintf("\"%v\"", val)
}

func InTags(tags ...string) string {
	return "{" + strings.Join(tags, "|") + "}"
}
