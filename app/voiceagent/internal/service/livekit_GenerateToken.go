package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"store/api/voiceagent"
	"store/confs"
	"store/pkg/krathelper"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// 1. Fetch User Profile
	var userNickname = "User"
	// Use GetById from core.go
	userProfile, err := s.data.Mongo.UserProfile.GetById(ctx, userId)
	if err == nil && userProfile != nil {
		if userProfile.Nickname != "" {
			userNickname = userProfile.Nickname
		}
	}

	// 2. Fetch User Memories (Important ones)
	// Using ListAndCount with bson.M requires importing 'go.mongodb.org/mongo-driver/bson'
	memoriesList, _, err := s.data.Mongo.Memory.ListAndCount(ctx, bson.M{"userId": userId}, options.Find().SetLimit(5).SetSort(bson.M{"createdAt": -1}))
	var memoryTexts []string
	if err == nil {
		for _, m := range memoriesList {
			// Simple filter: only send reasonably important memories
			if m.Importance >= 5 {
				memoryTexts = append(memoryTexts, m.Content)
			}
			if len(memoryTexts) >= 5 {
				break
			}
		}
	}

	// 3. Create Conversation Record
	conversation := &voiceagent.Conversation{
		XId:           primitive.NewObjectID().Hex(),
		UserId:        userId,
		AgentId:       "aura_zh", // Default or from req
		Status:        "pending",
		CreatedAt:     time.Now().Unix(),
		LastMessageAt: time.Now().Unix(),
	}
	_, err = s.data.Mongo.Conversation.Insert(ctx, conversation) // Ignore error? Or log it?
	if err != nil {
		return nil, err
	}

	// 4. Construct Metadata
	agentConfig := map[string]interface{}{
		"conversationId": conversation.XId,
		"agentName":      "aura_zh", // Explicitly specify the role
		"userId":         userId,
		"userProfile": map[string]string{
			"nickname": userNickname,
		},
		"memories": memoryTexts,
	}
	metadataJSON, _ := json.Marshal(agentConfig)

	// Explicitly create room with Agent config
	// "aura_zh" acts as the service pipe. The actual persona is driven by the Metadata inside.
	_, err = s.data.RoomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:            roomName,
		EmptyTimeout:    10 * 60,
		MaxParticipants: 2,
		Metadata:        string(metadataJSON),
		//Agents: []*livekit.RoomAgentDispatch{
		//	{
		//		AgentName: "aura_zh",
		//		Metadata:  string(metadataJSON),
		//	},
		//},
	})
	if err != nil {
		// Log error but proceed, as the token is still valid and room might be created on join
		// or already exists.
		fmt.Printf("Warning: Failed to create room explicitely (likely exists): %v\n", err)
		// Do not return error here, proceed to generate token.
	}

	return &voiceagent.GenerateLiveKitTokenResponse{
		AccessToken: token,
		Url:         confs.LiveKitUrl,
	}, nil
}
