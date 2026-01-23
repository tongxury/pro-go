package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type SettingsCollection struct {
	*mgz.Core[*projpb.AppSettings]
}

func NewSettingsCollection(db *mongo.Database) *SettingsCollection {
	return &SettingsCollection{
		Core: mgz.NewCore[*projpb.AppSettings](db, "settings"),
	}
}
