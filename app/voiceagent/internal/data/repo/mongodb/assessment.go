package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type AssessmentCollection struct {
	*mgz.Core[*voiceagent.Assessment]
}

func NewAssessmentCollection(db *mongo.Database) *AssessmentCollection {
	return &AssessmentCollection{
		Core: mgz.NewCore[*voiceagent.Assessment](db, "va_assessments"),
	}
}
