package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type AnswerCollection struct {
	*mgz.CoreCollection[aiagentpb.Answer]
}

func NewAnswerCollection(db *mongo.Database) *AnswerCollection {
	return &AnswerCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.Answer](db, "answers"),
	}
}
