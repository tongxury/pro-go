package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListConversations(ctx context.Context, req *voiceagent.ListConversationsRequest) (*voiceagent.ConversationList, error) {
	userId := krathelper.RequireUserId(ctx)

	list, total, err := s.Data.Mongo.Conversation.ListAndCount(ctx, bson.M{"user._id": userId}, mgz.Paging(req.Page, req.Size))
	if err != nil {
		return nil, err
	}

	return &voiceagent.ConversationList{
		List:  list,
		Total: total,
	}, nil
}
