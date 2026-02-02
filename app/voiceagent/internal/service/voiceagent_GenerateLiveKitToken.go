package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"store/api/voiceagent"
	"store/confs"

	"github.com/livekit/protocol/auth"
)

func (s *VoiceAgentService) GenerateLiveKitToken(ctx context.Context, req *voiceagent.GenerateLiveKitTokenRequest) (*voiceagent.GenerateLiveKitTokenResponse, error) {
	apiKey := confs.LiveKitApiKey
	apiSecret := confs.LiveKitApiSecret

	// 如果请求中没有指定房间名或身份，使用默认值
	roomName := req.RoomName
	if roomName == "" {
		roomName = "default_room"
	}

	identity := req.Identity
	if identity == "" {
		identity = "user_" + fmt.Sprintf("%d", time.Now().Unix())
	}

	at := auth.NewAccessToken(apiKey, apiSecret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomName,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	token, err := at.ToJWT()
	if err != nil {
		return nil, err
	}

	// 默认使用本地开发地址，可通过环境变量 LIVEKIT_URL 覆盖
	url := os.Getenv("LIVEKIT_URL")
	if url == "" {
		url = "ws://localhost:7880"
	}

	return &voiceagent.GenerateLiveKitTokenResponse{
		AccessToken: token,
		Url:         url,
	}, nil
}
