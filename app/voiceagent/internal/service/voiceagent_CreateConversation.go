package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *VoiceAgentService) CreateConversation(ctx context.Context, req *voiceagent.CreateConversationRequest) (*voiceagent.Conversation, error) {
	// Create a new conversation record
	userId := krathelper.RequireUserId(ctx)

	conv := &voiceagent.Conversation{
		XId:           primitive.NewObjectID().Hex(),
		UserId:        userId,
		AgentId:       req.AgentId,
		Status:        "active",
		CreatedAt:     time.Now().Unix(),
		LastMessageAt: time.Now().Unix(),
		// sceneId, etc. can be added if needed
	}

	_, err := s.Data.Mongo.Conversation.Insert(ctx, conv)
	if err != nil {
		return nil, err
	}

	return conv, nil
}
