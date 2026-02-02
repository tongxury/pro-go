package repo

import (
	demopb "store/api/demo"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

// FooCollection: Foo 资源的 MongoDB 仓储实现。
type FooCollection struct {
	*mgz.Core[*demopb.Foo]
}

func NewFooCollection(db *mongo.Database) *FooCollection {
	return &FooCollection{
		Core: mgz.NewCore[*demopb.Foo](db, "foos"),
	}
}
