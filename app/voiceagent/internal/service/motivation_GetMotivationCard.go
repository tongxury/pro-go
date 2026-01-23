package service

import (
	"context"
	voiceagent "store/api/voiceagent"
)

// GetMotivationCard: 获取激励卡片的详细信息。
func (s *VoiceAgentService) GetMotivationCard(ctx context.Context, req *voiceagent.GetMotivationCardRequest) (*voiceagent.MotivationCard, error) {
	card, err := s.Data.Mongo.Motivation.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	// 如果不是公开的，则需要验证权限（这里简单处理，如果是公开的谁都能看，否则只有本人能看）
	// userId, _ := krathelper.GetUserId(ctx)
	// if !card.IsPublic && card.UserId != userId {
	//     return nil, fmt.Errorf("unauthorized")
	// }
	return card, nil
}
