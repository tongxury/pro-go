package service

import (
	"context"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"time"
)

func (s *VoiceAgentService) AddTranscriptEntry(ctx context.Context, req *voiceagent.AddTranscriptEntryRequest) (*voiceagent.TranscriptEntry, error) {
	// TODO: Auth check if needed, but currently internal or Python Agent calls it.

	entry := &voiceagent.TranscriptEntry{
		User:         &ucpb.User{XId: req.UserId},
		Conversation: &voiceagent.Conversation{XId: req.ConversationId},
		Message:      req.Content,
		Role:         req.Role,
		CreatedAt:    time.Now().Unix(),
		AudioUrl:     req.AudioUrl,
	}

	_, err := s.Data.Mongo.Transcript.Insert(ctx, entry)
	if err != nil {
		return nil, err
	}

	// Update conversation status to active and refresh LastMessageAt
	if req.ConversationId != "" {
		// Use underlying collection to avoid defining generic update methods if not present
		_, err = s.Data.Mongo.Conversation.UpdateByIDIfExists(ctx, req.ConversationId,
			mgz.Op().Set("status", "ongoing").Set("lastMessageAt", time.Now().Unix()),
		)

		if err != nil {
			return nil, err
		}
	}

	return entry, nil
}
