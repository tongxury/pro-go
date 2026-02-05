package service

import (
	"context"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"

	"github.com/livekit/protocol/livekit"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *LiveKitService) StopConversation(ctx context.Context, req *voiceagent.StopConversationRequest) (*voiceagent.Conversation, error) {
	userId := krathelper.RequireUserId(ctx)

	// Verify conversation exists and belongs to user
	conv, err := s.data.Mongo.Conversation.FindByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "conversation not found")
	}

	if conv.User.XId != userId {
		return nil, status.Errorf(codes.PermissionDenied, "permission denied")
	}

	// Update status to completed
	updateOp := mgz.Op().
		Set("status", "completed").
		Set("endedAt", time.Now().Unix())

	_, err = s.data.Mongo.Conversation.UpdateByIDIfExists(ctx, req.Id, updateOp)
	if err != nil {
		return nil, err
	}

	// Clean up LiveKit Room
	if conv.RoomName != "" {
		_, err := s.data.RoomClient.DeleteRoom(ctx, &livekit.DeleteRoomRequest{
			Room: conv.RoomName,
		})
		if err != nil {
			fmt.Printf("Warning: Failed to delete room %s: %v\n", conv.RoomName, err)
		}
	}

	// Refetch to get the updated status
	updatedConv, err := s.data.Mongo.Conversation.FindByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return updatedConv, nil
}
