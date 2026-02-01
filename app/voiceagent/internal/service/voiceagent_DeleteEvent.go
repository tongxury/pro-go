package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *VoiceAgentService) DeleteEvent(ctx context.Context, req *voiceagent.DeleteEventRequest) (*emptypb.Empty, error) {
	userId := krathelper.RequireUserId(ctx)

	// 确保只能删除自己的事件
	event, err := s.Data.Mongo.ImportantEvent.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if event.UserId != userId {
		return nil, krathelper.ErrForbidden
	}

	err = s.Data.Mongo.ImportantEvent.DeleteByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
