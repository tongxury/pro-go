package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type AssetCollection struct {
	*mgz.Core[*projpb.Asset]
}

func NewAssetCollection(db *mongo.Database) *AssetCollection {
	return &AssetCollection{
		Core: mgz.NewCore[*projpb.Asset](db, "assets"),
	}
}
