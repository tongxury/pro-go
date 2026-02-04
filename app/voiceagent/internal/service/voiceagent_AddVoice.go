package service

import (
	"context"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AddVoice: 用户自定义克隆声音或系统预设。
func (s *VoiceAgentService) AddVoice(ctx context.Context, req *voiceagent.AddVoiceRequest) (*voiceagent.Voice, error) {
	userId := krathelper.RequireUserId(ctx)

	// 1. 调用 ElevenLabs SDK 进行声音克隆
	// 注意：此处需要处理文件下载或直接传入 URL
	// 为了演示，我们假设 AddVoice 接收一个采样 URL 并调用 SDK
	voiceId, err := s.Data.ElevenLabs.AddVoice(ctx, req.Name, "User cloned voice", []string{req.SampleUrl})
	if err != nil {
		return nil, err
	}

	// 2. 构建本地 Voice 对象
	v := &voiceagent.Voice{
		XId:       primitive.NewObjectID().Hex(),
		User:      &ucpb.User{XId: userId},
		Name:      req.Name,
		Type:      req.Type,
		VoiceId:   voiceId,
		Status:    "processing", // 克隆通常是异步的
		SampleUrl: req.SampleUrl,
		CreatedAt: time.Now().Unix(),
	}

	res, err := s.Data.Mongo.Voice.Insert(ctx, v)
	if err != nil {
		return nil, err
	}
	return res, nil
}
