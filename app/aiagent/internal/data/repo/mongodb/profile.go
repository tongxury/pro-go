package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type ProfileCollection struct {
	*mgz.CoreCollection[aiagentpb.Profile]
}

func NewProfileCollection(db *mongo.Database) *ProfileCollection {
	return &ProfileCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.Profile](db, "profiles"),
	}
}
