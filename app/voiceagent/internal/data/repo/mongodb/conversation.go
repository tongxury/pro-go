package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type ConversationCollection struct {
	*mgz.Core[*voiceagent.Conversation]
}

func NewConversationCollection(db *mongo.Database) *ConversationCollection {
	return &ConversationCollection{
		Core: mgz.NewCore[*voiceagent.Conversation](db, "va_conversations"),
	}
}
