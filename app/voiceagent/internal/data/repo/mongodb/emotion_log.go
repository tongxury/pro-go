package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type EmotionLogCollection struct {
	*mgz.Core[*voiceagent.EmotionLog]
}

func NewEmotionLogCollection(db *mongo.Database) *EmotionLogCollection {
	return &EmotionLogCollection{
		Core: mgz.NewCore[*voiceagent.EmotionLog](db, "va_emotion_logs"),
	}
}
