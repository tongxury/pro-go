package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *VoiceAgentService) SendMessage(ctx context.Context, req *voiceagent.SendMessageRequest) (*voiceagent.TranscriptEntry, error) {
	userId := krathelper.RequireUserId(ctx)

	sess, err := s.Data.Mongo.Conversation.GetById(ctx, req.ConversationId)
	if err != nil {
		return nil, err
	}

	userMsg := &voiceagent.TranscriptEntry{
		XId:            primitive.NewObjectID().Hex(),
		UserId:         userId,
		ConversationId: req.ConversationId,
		Message:        req.Message,
		Role:           "user",
		CreatedAt:      time.Now().Unix(),
	}
	if _, err := s.Data.Mongo.Transcript.Insert(ctx, userMsg); err != nil {
		return nil, err
	}

	aiText := "I'm processing your request via ElevenLabs..."
	aiAudioUrl := ""
	messageId := ""

	if req.EnableVoice {
		agent, _ := s.Data.Mongo.Agent.GetById(ctx, sess.AgentId)
		voiceId := agent.VoiceId
		if voiceId == "" {
			voiceId = "21m00Tcm4TlvDq8ikWAM"
		}
		aiAudioUrl = "https://api.elevenlabs.io/v1/text-to-speech/" + voiceId + "/stream"
	}

	aiMsg := &voiceagent.TranscriptEntry{
		XId:            primitive.NewObjectID().Hex(),
		UserId:         userId,
		ConversationId: req.ConversationId,
		Message:        aiText,
		AudioUrl:       aiAudioUrl,
		Role:           "agent",
		CreatedAt:      time.Now().Unix(),
		MessageId:      messageId,
	}
	res, aiErr := s.Data.Mongo.Transcript.Insert(ctx, aiMsg)
	if aiErr != nil {
		return nil, aiErr
	}

	updateOp := mgz.Op().Set("lastMessageAt", time.Now().Unix())
	_, _ = s.Data.Mongo.Conversation.UpdateByIDIfExists(ctx, sess.XId, updateOp)

	return res, nil
}
