package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type ItemCollection struct {
	*mgz.CoreCollection[aiagentpb.Item]
}

func NewItemCollection(db *mongo.Database) *ItemCollection {
	return &ItemCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.Item](db, "items"),
	}
}
