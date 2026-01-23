package redizv2

import (
	"fmt"
	"strings"
)

type Expr interface {
	String() string
}

type NumericRangeExpr[T any] struct {
	Start, End T
}

func (t NumericRangeExpr[T]) String() string {
	return fmt.Sprintf("[%v, %v]", t.Start, t.End)
}

type NumericLTExpr[T any] struct {
	Value T
}

func (t NumericLTExpr[T]) String() string {
	return fmt.Sprintf("[-inf, %v]", t.Value)
}

type NumericGTExpr[T any] struct {
	Value T
}

func (t NumericGTExpr[T]) String() string {
	return fmt.Sprintf("[%v, +inf]", t.Value)
}

// TagInExpr
type TagInExpr struct {
	Values []string
}

func (t TagInExpr) String() string {
	return "{" + strings.Join(t.Values, "|") + "}"
}

//

type TextContainsExpr struct {
	Value string
}

func (t TextContainsExpr) String() string {
	return "%" + t.Value + "%"
}

// RawExpr -
type RawExpr struct {
	Expr string
}

func (t RawExpr) String() string {
	return t.Expr
}

type Q interface {
	String() string
}

type Query struct {
	Field     string
	ValueExpr string
}

type AndQueries []Query

func (qs AndQueries) String() string {

	if len(qs) == 0 {
		return "*"
	}

	var queries []string
	for i := range qs {

		x := qs[i]

		queries = append(queries, fmt.Sprintf("@%s: %v", x.Field, x.ValueExpr))
	}

	return strings.Join(queries, " ")
}
