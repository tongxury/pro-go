package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListScenes(ctx context.Context, req *voiceagent.ListScenesRequest) (*voiceagent.SceneList, error) {
	// 场景通常为系统预设
	list, err := s.Data.Mongo.Scene.List(ctx, bson.M{}, mgz.Paging(0, 100))
	if err != nil {
		return nil, err
	}

	return &voiceagent.SceneList{
		List: list,
	}, nil
}
