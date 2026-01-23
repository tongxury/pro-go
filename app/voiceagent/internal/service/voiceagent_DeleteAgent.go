package service

import (
	"context"
	voiceagent "store/api/voiceagent"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *VoiceAgentService) DeleteAgent(ctx context.Context, req *voiceagent.DeleteAgentRequest) (*emptypb.Empty, error) {
	agent, err := s.Data.Mongo.Agent.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return &emptypb.Empty{}, nil
	}

	// 删除 ElevenLabs 端的 Agent
	if agent.AgentId != "" {
		_ = s.Data.ElevenLabs.DeleteAgent(ctx, agent.AgentId)
	}

	// 删除本地数据库记录
	err = s.Data.Mongo.Agent.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
