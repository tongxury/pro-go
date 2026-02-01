package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *VoiceAgentService) DeleteMemory(ctx context.Context, req *voiceagent.DeleteMemoryRequest) (*emptypb.Empty, error) {
	userId := krathelper.RequireUserId(ctx)

	// 确保只能删除自己的记忆
	memory, err := s.Data.Mongo.Memory.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if memory.UserId != userId {
		return nil, krathelper.ErrForbidden
	}

	err = s.Data.Mongo.Memory.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
