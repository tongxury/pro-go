package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListVoices(ctx context.Context, req *voiceagent.ListVoicesRequest) (*voiceagent.VoiceList, error) {
	userId := krathelper.RequireUserId(ctx)

	// 列出系统预设、用户克隆或全部声音
	var filter bson.M
	// 注意：数据库中存储的是 nested user object, key 为 "user._id"
	if req.Owner == "system" {
		filter = bson.M{"user._id": "system"}
	} else if req.Owner == "custom" {
		filter = bson.M{"user._id": userId}
	} else {
		// 默认返回 System + Custom
		filter = bson.M{
			"$or": []bson.M{
				{"user._id": userId},
				{"user._id": "system"},
			},
		}
	}

	list, err := s.Data.Mongo.Voice.List(ctx, filter, mgz.Paging(0, 100))
	if err != nil {
		return nil, err
	}

	return &voiceagent.VoiceList{
		List: list,
	}, nil
}
