package mgz

import (
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FindOptions struct {
	c *options.FindOptions
}

func Find() *FindOptions {
	return &FindOptions{
		c: options.Find(),
	}
}

func (t *FindOptions) B() *options.FindOptions {
	return t.c
}

func (t *FindOptions) Paging(page, size int64) *FindOptions {
	if size == 0 {
		return t
	}

	if page == 0 {
		t.c.SetLimit(int64(size))
		return t
	}

	t.c.SetLimit(int64(size)).SetSkip(int64((page - 1) * size))

	return t
}

func (t *FindOptions) SetFields(fields string) *FindOptions {

	if fields == "" {
		return t
	}

	fs := strings.Split(fields, ",")

	if len(fs) == 0 {
		return t
	}

	pro := bson.M{}
	for _, f := range fs {
		pro[f] = 1
	}

	t.c.SetProjection(pro)

	return t
}

func (t *FindOptions) SetSort(field string, d int) *FindOptions {

	t.c.SetSort(bson.M{field: d})

	return t
}

func (t *FindOptions) Projection(m bson.M) *FindOptions {
	t.c.SetProjection(m)
	return t
}

func (t *FindOptions) Limit(d int64) *FindOptions {

	t.c.SetLimit(d)

	return t
}

func (t *FindOptions) PageSize(p, s int64) *FindOptions {

	if s == 0 {
		return t
	}

	if p == 0 {
		t.c.SetLimit(s)
		return t
	}

	t.c.SetLimit(s).SetSkip((p - 1) * s)
	return t
}
