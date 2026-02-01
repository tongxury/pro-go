package repo

import (
	demopb "store/api/demo"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

// BarCollection: Bar 资源的 MongoDB 仓储实现。
type BarCollection struct {
	*mgz.Core[*demopb.Bar]
}

func NewBarCollection(db *mongo.Database) *BarCollection {
	return &BarCollection{
		Core: mgz.NewCore[*demopb.Bar](db, "bars"),
	}
}
