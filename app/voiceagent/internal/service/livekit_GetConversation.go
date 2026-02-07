package service

import (
	"context"
	voiceagent "store/api/voiceagent"
)

func (s *LiveKitService) GetConversation(ctx context.Context, req *voiceagent.GetConversationRequest) (*voiceagent.Conversation, error) {
	conv, err := s.data.Mongo.Conversation.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	if conv.Agent != nil && conv.Agent.XId != "" {
		agent, err := s.data.Mongo.Agent.GetById(ctx, conv.Agent.XId)
		if err == nil {
			conv.Agent = agent
		}
	}

	return conv, nil
}
