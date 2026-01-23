package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type TranscriptCollection struct {
	*mgz.Core[*voiceagent.TranscriptEntry]
}

func NewTranscriptCollection(db *mongo.Database) *TranscriptCollection {
	return &TranscriptCollection{
		Core: mgz.NewCore[*voiceagent.TranscriptEntry](db, "va_transcripts"),
	}
}
