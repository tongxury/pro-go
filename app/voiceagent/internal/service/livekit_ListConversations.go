package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *LiveKitService) ListConversations(ctx context.Context, req *voiceagent.ListConversationsRequest) (*voiceagent.ConversationList, error) {
	userId := krathelper.RequireUserId(ctx)

	list, total, err := s.data.Mongo.Conversation.ListAndCount(ctx,
		bson.M{"user._id": userId},
		mgz.Find().Paging(req.Page, req.Size).SetSort("createdAt", -1).B())
	if err != nil {
		return nil, err
	}

	return &voiceagent.ConversationList{
		List:  list,
		Total: total,
	}, nil
}
