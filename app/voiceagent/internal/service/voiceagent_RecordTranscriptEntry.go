package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *VoiceAgentService) RecordTranscriptEntry(ctx context.Context, req *voiceagent.RecordTranscriptEntryRequest) (*voiceagent.TranscriptEntry, error) {
	userId := krathelper.RequireUserId(ctx)

	entry := &voiceagent.TranscriptEntry{
		XId:            primitive.NewObjectID().Hex(),
		UserId:         userId,
		ConversationId: req.ConversationId,
		Role:           req.Role,
		Message:        req.Message,
		CreatedAt:      time.Now().Unix(),
	}

	res, err := s.Data.Mongo.Transcript.Insert(ctx, entry)
	if err != nil {
		return nil, err
	}

	// 更新最后消息时间
	updateOp := mgz.Op().Set("lastMessageAt", time.Now().Unix())
	_, _ = s.Data.Mongo.Conversation.UpdateByIDIfExists(ctx, req.ConversationId, updateOp)

	return res, nil
}
