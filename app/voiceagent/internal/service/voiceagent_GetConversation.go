package service

import (
	"context"
	voiceagent "store/api/voiceagent"
)

func (s *VoiceAgentService) GetConversation(ctx context.Context, req *voiceagent.GetConversationRequest) (*voiceagent.Conversation, error) {
	return s.Data.Mongo.Conversation.GetById(ctx, req.Id)
}
