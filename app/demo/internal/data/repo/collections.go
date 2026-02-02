package repo

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type Collections struct {
	Foo *FooCollection
	Bar *BarCollection
}

func NewCollections(db *mongo.Database) *Collections {
	return &Collections{
		Foo: NewFooCollection(db),
		Bar: NewBarCollection(db),
	}
}
