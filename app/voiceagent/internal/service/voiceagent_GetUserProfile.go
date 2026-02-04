package service

import (
	"context"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *VoiceAgentService) GetUserProfile(ctx context.Context, req *voiceagent.GetUserProfileRequest) (*voiceagent.UserProfile, error) {
	userId := krathelper.RequireUserId(ctx)

	profile, err := s.Data.Mongo.UserProfile.FindOne(ctx, bson.M{"user._id": userId})
	if err != nil {
		// 如果不存在，创建一个空的档案
		profile = &voiceagent.UserProfile{
			XId:       primitive.NewObjectID().Hex(),
			User:      &ucpb.User{XId: userId},
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}
		_, err = s.Data.Mongo.UserProfile.Insert(ctx, profile)
		if err != nil {
			return nil, err
		}
	}

	return profile, nil
}
