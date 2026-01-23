package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommodityCollection struct {
	*mgz.Core[*projpb.Commodity]
}

func NewCommodityCollection(db *mongo.Database) *CommodityCollection {
	return &CommodityCollection{
		Core: mgz.NewCore[*projpb.Commodity](db, "commodities"),
	}
}
