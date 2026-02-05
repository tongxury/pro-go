package service

import (
	"context"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"
)

func (s *LiveKitService) UpdateConversation(ctx context.Context, req *voiceagent.UpdateConversationRequest) (*voiceagent.Conversation, error) {

	userId := krathelper.RequireUserId(ctx)

	// 1. Get conversation to verify ownership and existence
	conv, err := s.data.Mongo.Conversation.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if conv == nil {
		return nil, fmt.Errorf("conversation not found")
	}
	if conv.User.GetXId() != userId {
		return nil, krathelper.ErrForbidden
	}

	switch req.Action {
	case "end":

		updateOp := mgz.Op().
			Set("status", "ended").
			Set("endedAt", time.Now().Unix())

		_, err := s.data.Mongo.Conversation.UpdateByIDIfExists(ctx, req.Id, updateOp)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}
