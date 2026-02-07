package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"

	mapset "github.com/deckarep/golang-set/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *LiveKitService) ListConversations(ctx context.Context, req *voiceagent.ListConversationsRequest) (*voiceagent.ConversationList, error) {
	userId := krathelper.RequireUserId(ctx)

	list, total, err := s.data.Mongo.Conversation.ListAndCount(ctx,
		bson.M{"user._id": userId, "status": bson.M{"$in": []string{"completed"}}},
		mgz.Find().Paging(req.Page, req.Size).SetSort("createdAt", -1).B())
	if err != nil {
		return nil, err
	}

	// 补充 agent 信息
	if len(list) > 0 {
		agentIdsSet := mapset.NewSet[string]()
		for _, conv := range list {
			if conv.Agent != nil && conv.Agent.XId != "" {
				agentIdsSet.Add(conv.Agent.XId)
			}
		}

		agentIds := agentIdsSet.ToSlice()

		if len(agentIds) > 0 {
			agents, err := s.data.Mongo.Agent.List(ctx, mgz.Filter().Ids(agentIds).B())
			if err == nil {
				agentMap := make(map[string]*voiceagent.Agent)
				for _, a := range agents {
					agentMap[a.XId] = a
				}

				for _, conv := range list {
					if conv.Agent != nil {
						if agent, ok := agentMap[conv.Agent.XId]; ok {
							conv.Agent = agent
						}
					}
				}
			}
		}
	}

	return &voiceagent.ConversationList{
		List:  list,
		Total: total,
	}, nil
}
