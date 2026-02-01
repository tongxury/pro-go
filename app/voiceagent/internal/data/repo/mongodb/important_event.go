package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type ImportantEventCollection struct {
	*mgz.Core[*voiceagent.ImportantEvent]
}

func NewImportantEventCollection(db *mongo.Database) *ImportantEventCollection {
	return &ImportantEventCollection{
		Core: mgz.NewCore[*voiceagent.ImportantEvent](db, "va_important_events"),
	}
}
