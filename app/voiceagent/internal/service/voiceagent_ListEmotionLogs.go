package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListEmotionLogs(ctx context.Context, req *voiceagent.ListEmotionLogsRequest) (*voiceagent.EmotionLogList, error) {
	userId := krathelper.RequireUserId(ctx)

	filter := bson.M{"userId": userId}

	// 添加时间范围过滤
	if req.StartTime > 0 || req.EndTime > 0 {
		timeFilter := bson.M{}
		if req.StartTime > 0 {
			timeFilter["$gte"] = req.StartTime
		}
		if req.EndTime > 0 {
			timeFilter["$lte"] = req.EndTime
		}
		filter["createdAt"] = timeFilter
	}

	page := req.Page
	size := req.Size
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 30
	}

	list, total, err := s.Data.Mongo.EmotionLog.ListWithFilterAndSort(ctx, filter, bson.M{"createdAt": -1}, page, size)
	if err != nil {
		return nil, err
	}

	return &voiceagent.EmotionLogList{
		List:  list,
		Total: total,
	}, nil
}
