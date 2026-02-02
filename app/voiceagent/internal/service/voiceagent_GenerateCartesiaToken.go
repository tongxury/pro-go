package service

import (
	"context"
	voiceagent "store/api/voiceagent"
)

func (s *VoiceAgentService) GenerateCartesiaToken(ctx context.Context, req *voiceagent.GenerateCartesiaTokenRequest) (*voiceagent.GenerateCartesiaTokenResponse, error) {
	// 为 TTS 申请权限，有效期 1 小时 (3600秒)
	token, err := s.Data.Cartesia.GenerateAccessToken(ctx, map[string]bool{"tts": true}, 3600)
	if err != nil {
		return nil, err
	}

	return &voiceagent.GenerateCartesiaTokenResponse{
		AccessToken: token,
	}, nil
}
