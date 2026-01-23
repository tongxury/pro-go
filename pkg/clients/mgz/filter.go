package mgz

import (
	"go.mongodb.org/mongo-driver/bson"
)

type FilterOptions struct {
	c bson.M
}

func Filter() *FilterOptions {
	return &FilterOptions{
		c: bson.M{},
	}
}

func (t *FilterOptions) B() bson.M {
	return t.c
}

func (t *FilterOptions) Ids(values []string) *FilterOptions {

	var ids []any
	for _, v := range values {
		ids = append(ids, ObjectId(v))
	}

	return t.In("_id", ids)
}

func (t *FilterOptions) HasPrefix(field string, value string) *FilterOptions {
	t.c[field] = bson.M{"$regex": "^" + value, "$options": "i"}
	return t
}

func (t *FilterOptions) InString(field string, values []string) *FilterOptions {

	var vs []interface{}
	for _, v := range values {
		if v != "" {
			vs = append(vs, v)
		}
	}

	if len(vs) == 0 {
		return t
	}

	t.c[field] = bson.M{
		"$in": vs,
	}

	return t
}

func (t *FilterOptions) In(field string, values []any) *FilterOptions {

	var vs []interface{}
	for _, v := range values {
		if v != "" {
			vs = append(vs, v)
		}
	}

	if len(vs) == 0 {
		return t
	}

	t.c[field] = bson.M{
		"$in": vs,
	}

	return t
}

func (t *FilterOptions) EQ(field string, value any) *FilterOptions {
	if value == "" {
		return t
	}

	t.c[field] = value

	return t
}

func (t *FilterOptions) GT(field string, value any) *FilterOptions {
	if value == "" {
		return t
	}

	t.c[field] = bson.M{"$gt": value}
	return t
}

func (t *FilterOptions) LT(field string, value any) *FilterOptions {
	if value == "" {
		return t
	}

	t.c[field] = bson.M{"$lt": value}
	return t
}

func (t *FilterOptions) Exists(field string) *FilterOptions {
	if field == "" {
		return t
	}

	t.c[field] = bson.M{"$exists": true}

	return t
}

func (t *FilterOptions) NotExists(field string) *FilterOptions {
	if field == "" {
		return t
	}

	t.c[field] = bson.M{"$exists": false}

	return t
}

func (t *FilterOptions) NotEQ(field string, value any) *FilterOptions {
	if value == "" {
		return t
	}

	t.c[field] = bson.M{"$ne": value}

	return t
}

func (t *FilterOptions) And(filters ...*FilterOptions) *FilterOptions {
	if len(filters) == 0 {
		return t
	}

	var andList []interface{}
	if existing, ok := t.c["$and"]; ok {
		if list, ok := existing.([]interface{}); ok {
			andList = list
		}
	}

	for _, f := range filters {
		if len(f.c) > 0 {
			andList = append(andList, f.c)
		}
	}

	if len(andList) > 0 {
		t.c["$and"] = andList
	}

	return t
}
