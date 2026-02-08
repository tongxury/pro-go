package service

import (
	"context"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) UpdateUserProfile(ctx context.Context, req *voiceagent.UpdateUserProfileRequest) (*voiceagent.UserProfile, error) {
	userId := krathelper.RequireUserId(ctx)

	// 先获取或创建档案
	profile, err := s.Data.Mongo.UserProfile.FindOne(ctx, bson.M{"user._id": userId})
	if err != nil || profile == nil {
		// 如果不存在，创建新档案
		profile = &voiceagent.UserProfile{
			User:        &ucpb.User{XId: userId},
			Nickname:    req.Nickname,
			Birthday:    req.Birthday,
			Interests:   req.Interests,
			Goals:       req.Goals,
			Bio:         req.Bio,
			Personality: req.Personality,
			Timezone:    req.Timezone,
			CreatedAt:   time.Now().Unix(),
			UpdatedAt:   time.Now().Unix(),
		}
		res, err := s.Data.Mongo.UserProfile.Insert(ctx, profile)
		if err != nil {
			return nil, err
		}
		return res, nil
	}

	// 构建更新操作
	updateOp := mgz.Op().Set("updatedAt", time.Now().Unix())

	if req.Nickname != "" {
		updateOp = updateOp.Set("nickname", req.Nickname)
	}
	if req.Birthday != "" {
		updateOp = updateOp.Set("birthday", req.Birthday)
	}
	if len(req.Interests) > 0 {
		updateOp = updateOp.Set("interests", req.Interests)
	}
	if len(req.Goals) > 0 {
		updateOp = updateOp.Set("goals", req.Goals)
	}
	if req.Bio != "" {
		updateOp = updateOp.Set("bio", req.Bio)
	}
	if req.Personality != "" {
		updateOp = updateOp.Set("personality", req.Personality)
	}
	if req.Timezone != "" {
		updateOp = updateOp.Set("timezone", req.Timezone)
	}

	_, err = s.Data.Mongo.UserProfile.UpdateByIDIfExists(ctx, profile.XId, updateOp)
	if err != nil {
		return nil, err
	}

	// 返回更新后的档案
	return s.Data.Mongo.UserProfile.GetById(ctx, profile.XId)
}
