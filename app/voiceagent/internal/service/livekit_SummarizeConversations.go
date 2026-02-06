package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
)

func (s *LiveKitService) SummarizeConversations(ctx context.Context, _ *empty.Empty) (*empty.Empty, error) {
	err := s.agentBiz.SummarizeEndedConversations(ctx)
	if err != nil {
		return nil, err
	}
	return &empty.Empty{}, nil
}
