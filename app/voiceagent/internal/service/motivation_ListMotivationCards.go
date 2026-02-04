package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"

	"go.mongodb.org/mongo-driver/bson"
)

// ListMotivationCards: 分页列出当前用户创建的所有激励卡片。
func (s *VoiceAgentService) ListMotivationCards(ctx context.Context, req *voiceagent.ListMotivationCardsRequest) (*voiceagent.MotivationCardList, error) {
	userId := krathelper.RequireUserId(ctx)

	filter := bson.M{"user._id": userId}
	list, total, err := s.Data.Mongo.Motivation.ListAndCount(ctx, filter, mgz.Paging(req.Page, req.Size))
	if err != nil {
		return nil, err
	}

	return &voiceagent.MotivationCardList{
		List:  list,
		Total: total,
	}, nil
}
