package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type SurveyCollection struct {
	*mgz.CoreCollection[aiagentpb.Survey]
}

func NewSurveyCollection(db *mongo.Database) *SurveyCollection {
	return &SurveyCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.Survey](db, "surveys"),
	}
}
