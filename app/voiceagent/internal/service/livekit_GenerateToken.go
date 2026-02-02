package service

import (
	"context"
	"fmt"
	"time"

	"store/api/voiceagent"
	"store/confs"
	"store/pkg/krathelper"

	"github.com/livekit/protocol/auth"
)

func (s *LiveKitService) GenerateLiveKitToken(ctx context.Context, req *voiceagent.GenerateLiveKitTokenRequest) (*voiceagent.GenerateLiveKitTokenResponse, error) {
	userId := krathelper.RequireUserId(ctx)

	apiKey := confs.LiveKitApiKey
	apiSecret := confs.LiveKitApiSecret

	// 如果请求中没有指定房间名或身份，使用默认值
	roomName := req.RoomName
	if roomName == "" {
		roomName = "default_room"
	}

	identity := req.Identity
	if identity == "" {
		identity = fmt.Sprintf("user_%s_%d", userId, time.Now().Unix())
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

	return &voiceagent.GenerateLiveKitTokenResponse{
		AccessToken: token,
		Url:         confs.LiveKitUrl,
	}, nil
}
