package service

import (
	"context"
	voiceagent "store/api/voiceagent"
)

func (s *LiveKitService) GetConversation(ctx context.Context, req *voiceagent.GetConversationRequest) (*voiceagent.Conversation, error) {
	return s.data.Mongo.Conversation.GetById(ctx, req.Id)
}
