package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListAgents(ctx context.Context, req *voiceagent.ListAgentsRequest) (*voiceagent.AgentList, error) {
	userId := krathelper.RequireUserId(ctx)

	filter := bson.M{"userId": userId}
	if req.Category != "" {
		filter["category"] = req.Category
	}

	list, total, err := s.Data.Mongo.Agent.ListAndCount(ctx, filter, mgz.Paging(req.Page, req.Size))
	if err != nil {
		return nil, err
	}

	return &voiceagent.AgentList{
		List:  list,
		Total: total,
	}, nil
}
