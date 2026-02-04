package service

import (
	"context"
	"fmt"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/pkg/krathelper"
	"store/pkg/sdk/third/elevenlabs"
	"store/pkg/sdk/third/gemini"
	"time"

	"google.golang.org/genai"
)

// CreateMotivationCard: 生成一段带情感的语音激励卡片。
func (s *VoiceAgentService) CreateMotivationCard(ctx context.Context, req *voiceagent.CreateMotivationCardRequest) (*voiceagent.MotivationCardResponse, error) {
	userId := krathelper.RequireUserId(ctx)

	// 1. 获取音色/助理信息
	agent, err := s.Data.Mongo.Agent.GetById(ctx, req.AgentId)
	if err != nil {
		return nil, err
	}
	if agent == nil {
		return nil, fmt.Errorf("agent not found")
	}

	// 1.5 获取用户信息 (用于海报展示)
	var userName, userAvatar string
	if s.Data.GrpcClients.UserCenterClient != nil {
		userResp, err := s.Data.GrpcClients.UserCenterClient.GetUser(ctx, &ucpb.GetUserRequest{})
		if err == nil && userResp != nil {
			userName = userResp.Nickname
			if userName == "" {
				userName = userResp.Name
			}
			userAvatar = userResp.Avatar
		}
	}

	// 2. 文案润色 (接入 Gemini)
	prompt := fmt.Sprintf("你是一个情感治愈与励志大师。请根据用户提供的草稿文案，结合'%s'的情感基调，将其润色成一段富有情感、优美且动人的短语。要求：1. 长度在30-80字之间；2. 适合用第一人称朗读；3. 不要包含任何特殊符号；4. 直接返回润色后的内容，不要有任何开场白。\n\n用户草稿：%s", req.EmotionTag, req.OriginalText)
	polishedText, err := s.Data.Gemini.Get().GenerateContent(ctx, gemini.GenerateContentRequest{
		Parts: []*genai.Part{{Text: prompt}},
	})
	if err != nil {
		// 如果 LLM 失败，回退到原始文案
		polishedText = req.OriginalText
	}

	// 3. 调用 ElevenLabs TTS 生成音频
	// modelId := req.ModelId
	// if modelId == "" {
	// }
	modelId := "eleven_multilingual_v2"

	voiceSettings := &elevenlabs.VoiceSettings{
		Stability:       0.5,
		SimilarityBoost: 0.75,
		Style:           0.0,
		UseSpeakerBoost: true,
	}

	audioData, err := s.Data.ElevenLabs.TextToSpeech(ctx, agent.VoiceId, &elevenlabs.TextToSpeechRequest{
		Text:          polishedText,
		ModelID:       modelId,
		VoiceSettings: voiceSettings,
	})
	if err != nil {
		return nil, err
	}

	// 4. 上传音频到 OSS
	audioUrl, err := s.Data.TOS.PutAudioBytes(ctx, audioData)
	if err != nil {
		return nil, err
	}

	// 4.5 生成波形数据 (简单模拟，用于前端声纹海报展示)
	// 实际场景可使用 ffmpeg 或 go 语音库提取能量值
	waveform := make([]float32, 40)
	for i := 0; i < 40; i++ {
		// 生成 0.1 到 1.0 之间的随机波动值，模拟声纹起伏
		waveform[i] = 0.1 + (float32(i%7) * 0.12) + (float32(i%3) * 0.05)
		if waveform[i] > 1.0 {
			waveform[i] = 0.9
		}
	}

	// 5. 保存卡片到本地数据库
	card := &voiceagent.MotivationCard{
		User: &ucpb.User{
			XId:    userId,
			Name:   userName,
			Avatar: userAvatar,
		},
		Agent:       agent,
		Text:        polishedText,
		AudioUrl:    audioUrl,
		EmotionTag:  req.EmotionTag,
		CreatedAt:   time.Now().Unix(),
		IsPublic:    req.IsPublic,
		Waveform:    waveform,
		PosterStyle: req.PosterStyle,
	}

	res, err := s.Data.Mongo.Motivation.Insert(ctx, card)
	if err != nil {
		return nil, err
	}

	// 6. 构建分享链接 (指向 H5 预览页)
	shareUrl := fmt.Sprintf("https://voiceagent.ai/share/%s", res.XId)
	qrCodeUrl := fmt.Sprintf("https://api.qrserver.com/v1/create-qr-code/?size=200x200&data=%s", shareUrl)

	// 更新数据库中的二维码链接
	res.QrCodeUrl = qrCodeUrl
	_, _ = s.Data.Mongo.Motivation.ReplaceByID(ctx, res.XId, res)

	return &voiceagent.MotivationCardResponse{
		Id:           res.XId,
		AudioUrl:     res.AudioUrl,
		ShareUrl:     shareUrl,
		PolishedText: res.Text,
		User:         res.User,
		Agent:        res.Agent,
		EmotionTag:   req.EmotionTag,
		CreatedAt:    res.CreatedAt,
		Waveform:     res.Waveform,
		PosterStyle:  res.PosterStyle,
		PosterUrl:    res.PosterUrl,
		QrCodeUrl:    qrCodeUrl,
	}, nil
}
