package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type VoiceCollection struct {
	*mgz.Core[*voiceagent.Voice]
}

func NewVoiceCollection(db *mongo.Database) *VoiceCollection {
	return &VoiceCollection{
		Core: mgz.NewCore[*voiceagent.Voice](db, "va_voices"),
	}
}
