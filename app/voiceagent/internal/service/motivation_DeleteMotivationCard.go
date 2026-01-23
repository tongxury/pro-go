package service

import (
	"context"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"

	"google.golang.org/protobuf/types/known/emptypb"
)

// DeleteMotivationCard: 删除特定的激励卡片。
func (s *VoiceAgentService) DeleteMotivationCard(ctx context.Context, req *voiceagent.DeleteMotivationCardRequest) (*emptypb.Empty, error) {
	userId := krathelper.RequireUserId(ctx)

	card, err := s.Data.Mongo.Motivation.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if card == nil || card.User.GetId() != userId {
		return nil, fmt.Errorf("motivation card not found or unauthorized")
	}

	err = s.Data.Mongo.Motivation.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
