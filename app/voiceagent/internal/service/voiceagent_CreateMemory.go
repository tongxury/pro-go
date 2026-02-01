package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *VoiceAgentService) CreateMemory(ctx context.Context, req *voiceagent.CreateMemoryRequest) (*voiceagent.Memory, error) {
	userId := krathelper.RequireUserId(ctx)

	memory := &voiceagent.Memory{
		XId:        primitive.NewObjectID().Hex(),
		UserId:     userId,
		Type:       req.Type,
		Content:    req.Content,
		Importance: req.Importance,
		Tags:       req.Tags,
		CreatedAt:  time.Now().Unix(),
		UpdatedAt:  time.Now().Unix(),
	}

	if memory.Importance == 0 {
		memory.Importance = 5 // 默认中等重要性
	}

	res, err := s.Data.Mongo.Memory.Insert(ctx, memory)
	if err != nil {
		return nil, err
	}

	return res, nil
}
