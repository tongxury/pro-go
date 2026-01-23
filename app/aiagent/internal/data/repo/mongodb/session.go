package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type SessionCollection struct {
	*mgz.CoreCollection[aiagentpb.Session]
}

func NewSessionCollection(db *mongo.Database) *SessionCollection {
	return &SessionCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.Session](db, "sessions"),
	}
}
