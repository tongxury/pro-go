package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *VoiceAgentService) CreateConversation(ctx context.Context, req *voiceagent.CreateConversationRequest) (*voiceagent.Conversation, error) {
	userId := krathelper.RequireUserId(ctx)

	agent, err := s.Data.Mongo.Agent.GetById(ctx, req.AgentId)
	if err != nil {
		return nil, err
	}

	if req.SceneId != "" {
		if _, err := s.Data.Mongo.Scene.GetById(ctx, req.SceneId); err != nil {
			return nil, err
		}
	}

	//signedURL := ""
	//if agent.AgentId != "" {
	//	signedURL, err = s.Data.ElevenLabs.GetSignedURL(ctx, agent.AgentId)
	//	if err != nil {
	//		return nil, err
	//	}
	//}

	token := ""
	if agent.AgentId != "" {
		token, err = s.Data.ElevenLabs.GenerateConversationToken(ctx, agent.AgentId)
		if err != nil {
			return nil, err
		}
	}

	sess := &voiceagent.Conversation{
		XId:       primitive.NewObjectID().Hex(),
		UserId:    userId,
		AgentId:   req.AgentId,
		SceneId:   req.SceneId,
		CreatedAt: time.Now().Unix(),
		Status:    "active",
		Token:     token,
		//SignedUrl: signedURL,
	}

	res, err := s.Data.Mongo.Conversation.Insert(ctx, sess)
	if err != nil {
		return nil, err
	}
	return res, nil
}
