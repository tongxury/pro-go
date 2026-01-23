package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
)

func (s *VoiceAgentService) UpdateConversation(ctx context.Context, req *voiceagent.UpdateConversationRequest) (*voiceagent.Conversation, error) {
	updateOp := mgz.Op()
	if req.Status != "" {
		updateOp.Set("status", req.Status)
	}
	if req.ConversationId != "" {
		updateOp.Set("conversationId", req.ConversationId)
	}

	_, err := s.Data.Mongo.Conversation.UpdateByIDIfExists(ctx, req.Id, updateOp)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
