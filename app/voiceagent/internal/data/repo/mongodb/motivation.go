package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"go.mongodb.org/mongo-driver/mongo"
)

type MotivationCollection struct {
	*mgz.Core[*voiceagent.MotivationCard]
}

func NewMotivationCollection(db *mongo.Database) *MotivationCollection {
	return &MotivationCollection{
		Core: mgz.NewCore[*voiceagent.MotivationCard](db, "va_motivations"),
	}
}
