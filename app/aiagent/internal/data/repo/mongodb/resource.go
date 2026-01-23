package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type ResourceCollection struct {
	*mgz.CoreCollection[aiagentpb.ResourceV2]
}

func NewResourceCollection(db *mongo.Database) *ResourceCollection {
	return &ResourceCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.ResourceV2](db, "resources"),
	}
}
