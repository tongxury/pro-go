package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListAgents(ctx context.Context, req *voiceagent.ListAgentsRequest) (*voiceagent.AgentList, error) {
	userId := krathelper.RequireUserId(ctx)

	filter := bson.M{"user._id": userId}
	if req.Category != "" {
		filter["category"] = req.Category
	}

	list, total, err := s.Data.Mongo.Agent.ListAndCount(ctx, filter, mgz.Paging(req.Page, req.Size))
	if err != nil {
		return nil, err
	}

	if total == 0 && req.Page <= 1 {
		err = s.AgentBiz.CreateDefaultCounselingAgent(ctx, userId)
		if err != nil {
			log.Errorw("failed to create default counseling agent in ListAgents", "error", err, "userId", userId)
		} else {
			// 再次查询
			list, total, err = s.Data.Mongo.Agent.ListAndCount(ctx, filter, mgz.Paging(req.Page, req.Size))
			if err != nil {
				return nil, err
			}
		}
	}

	return &voiceagent.AgentList{
		List:  list,
		Total: total,
	}, nil
}
