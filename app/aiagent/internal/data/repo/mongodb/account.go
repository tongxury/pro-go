package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type AccountCollection struct {
	*mgz.CoreCollectionV3[aiagentpb.Account]
}

func NewAccountCollection(db *mongo.Database) *AccountCollection {
	return &AccountCollection{
		CoreCollectionV3: mgz.NewCoreCollectionV3[aiagentpb.Account](db, "accounts"),
	}
}
