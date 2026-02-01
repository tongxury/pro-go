package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListMemories(ctx context.Context, req *voiceagent.ListMemoriesRequest) (*voiceagent.MemoryList, error) {
	userId := krathelper.RequireUserId(ctx)

	filter := bson.M{"userId": userId}
	if req.Type != "" {
		filter["type"] = req.Type
	}

	page := req.Page
	size := req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}

	list, total, err := s.Data.Mongo.Memory.ListWithFilterAndSort(ctx, filter, bson.M{"createdAt": -1}, page, size)
	if err != nil {
		return nil, err
	}

	return &voiceagent.MemoryList{
		List:  list,
		Total: total,
	}, nil
}
