package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	aiagentpb "store/api/aiagent"
	"store/app/aiagent/configs"
)

func (t TrackerService) GetPromptSettings(ctx context.Context, _ *emptypb.Empty) (*aiagentpb.PromptSettings, error) {
	return &aiagentpb.PromptSettings{
		Settings: configs.GetPromptSettings(),
	}, nil
}
