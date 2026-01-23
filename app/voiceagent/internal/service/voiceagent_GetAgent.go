package service

import (
	"context"
	voiceagent "store/api/voiceagent"
)

func (s *VoiceAgentService) GetAgent(ctx context.Context, req *voiceagent.GetAgentRequest) (*voiceagent.Agent, error) {
	return s.Data.Mongo.Agent.GetById(ctx, req.Id)
}
