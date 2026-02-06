package service

import (
	"context"
	"encoding/json"
	"fmt"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/confs"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *LiveKitService) CreateConversation(ctx context.Context, req *voiceagent.CreateConversationRequest) (*voiceagent.Conversation, error) {
	userId := krathelper.RequireUserId(ctx)

	// Generate RoomName and Identity
	roomName := fmt.Sprintf("room_%s_%d", userId, time.Now().Unix())

	// Create Conversation Record
	conv := &voiceagent.Conversation{
		XId:           primitive.NewObjectID().Hex(),
		User:          &ucpb.User{XId: userId},
		Agent:         &voiceagent.Agent{XId: req.AgentId},
		Status:        "pending",
		CreatedAt:     time.Now().Unix(),
		LastMessageAt: time.Now().Unix(),
		RoomName:      roomName,
	}

	_, err := s.data.Mongo.Conversation.Insert(ctx, conv)
	if err != nil {
		return nil, err
	}

	// Delegate join logic (context gathering, room creation, dispatch, token gen)
	token, url, err := s.joinRoom(ctx, userId, req.AgentId, conv.XId, roomName)
	if err != nil {
		// Note: Conversation record exists but room join failed.
		// We might want to log this or consider cleanup, but for now just return error.
		return nil, err
	}

	conv.Token = token
	conv.SignedUrl = url

	return conv, nil
}

func (s *LiveKitService) joinRoom(ctx context.Context, userId, agentId, conversationId, roomName string) (string, string, error) {
	// 1. Fetch Business Context: User Profile
	var userNickname = "User"
	userProfile, err := s.data.Mongo.UserProfile.FindOne(ctx, bson.M{"user._id": userId})
	if err == nil && userProfile != nil && userProfile.Nickname != "" {
		userNickname = userProfile.Nickname
	}

	// 2. Fetch Business Context: Memories
	memoriesList, _, err := s.data.Mongo.Memory.ListAndCount(ctx, bson.M{"user._id": userId}, options.Find().SetLimit(5).SetSort(bson.M{"createdAt": -1}))
	var memoryTexts []string
	if err == nil {
		for _, m := range memoriesList {
			if m.Importance >= 3 {
				memoryTexts = append(memoryTexts, conv.S2J(m))
			}
			if len(memoryTexts) >= 10 {
				break
			}
		}
	}

	// 3. Construct Business Metadata for the Agent
	agentConfig := map[string]interface{}{
		"conversationId": conversationId,
		"agentId":        agentId,
		"userId":         userId,
		"nickname":       userNickname,
		"memories":       memoryTexts,
		"agentName":      "aura_zh", // Base service pipe. TODO: Maybe fetch from Agent definition?
	}
	metadataJSON, _ := json.Marshal(agentConfig)

	log.Debugw("metadataJSON", string(metadataJSON))

	// 4. Create Room with Metadata
	_, err = s.data.RoomClient.CreateRoom(ctx, &livekit.CreateRoomRequest{
		Name:            roomName,
		EmptyTimeout:    10 * 60,
		MaxParticipants: 2,
		Metadata:        string(metadataJSON),
	})
	if err != nil {
		fmt.Printf("Warning: Failed to create room: %v\n", err)
	}

	// 5. Dispatch Agent
	_, err = s.data.AgentClient.CreateDispatch(ctx, &livekit.CreateAgentDispatchRequest{
		AgentName: "aura_zh",
		Room:      roomName,
		Metadata:  string(metadataJSON),
	})
	if err != nil {
		fmt.Printf("Warning: Failed to dispatch agent: %v\n", err)
	}

	// 6. Generate LiveKit Token
	apiKey := confs.LiveKitApiKey
	apiSecret := confs.LiveKitApiSecret
	identity := fmt.Sprintf("user_%s_%d", userId, time.Now().Unix())

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
		return "", "", err
	}

	return token, confs.LiveKitUrl, nil
}
