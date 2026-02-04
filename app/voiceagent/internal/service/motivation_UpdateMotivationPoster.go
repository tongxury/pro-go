package service

import (
	"context"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"

	"google.golang.org/protobuf/types/known/emptypb"
)

// UpdateMotivationPoster: 更新卡片的海报 URL。
func (s *VoiceAgentService) UpdateMotivationPoster(ctx context.Context, req *voiceagent.UpdateMotivationPosterRequest) (*emptypb.Empty, error) {
	userId := krathelper.RequireUserId(ctx)

	// 1. 获取卡片
	card, err := s.Data.Mongo.Motivation.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if card == nil {
		return nil, fmt.Errorf("motivation card not found")
	}

	// 2. 权限校验：只有创建者可以更新海报
	if card.GetUser().GetXId() != userId {
		return nil, fmt.Errorf("unauthorized to update this poster")
	}

	// 3. 更新海报 URL
	card.PosterUrl = req.PosterUrl
	_, err = s.Data.Mongo.Motivation.ReplaceByID(ctx, card.XId, card)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
