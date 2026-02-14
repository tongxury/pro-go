package service

import (
	"context"
	"fmt"
	voiceagent "store/api/voiceagent"
	"store/pkg/sdk/third/bytedance/tos"
	"time"
)

func (s *VoiceAgentService) PreviewVoice(ctx context.Context, req *voiceagent.PreviewVoiceRequest) (*voiceagent.PreviewVoiceResponse, error) {
	audioUrl, err := s.GeneratePreviewAudio(ctx, req.VoiceId, req.Text)
	if err != nil {
		return nil, err
	}

	return &voiceagent.PreviewVoiceResponse{
		AudioUrl: audioUrl,
	}, nil
}

func (s *VoiceAgentService) GeneratePreviewAudio(ctx context.Context, voiceId string, text string) (string, error) {
	if text == "" {
		// Default preview text in Chinese
		text = "你好！这是我为您准备的试听音色。希望你喜欢这个声音！"
	}

	// 1. Generate audio bytes from Cartesia
	audioBytes, err := s.Data.Cartesia.TextToSpeechBytes(ctx, voiceId, "sonic-multilingual", text)
	if err != nil {
		return "", err
	}

	return s.uploadPreview(ctx, voiceId, audioBytes)
}

func (s *VoiceAgentService) uploadPreview(ctx context.Context, voiceId string, audioBytes []byte) (string, error) {
	// 2. Upload to TOS
	// We use a unique key based on MD5 and timestamp to avoid collisions
	fileName := fmt.Sprintf("previews/%s/%d.wav", voiceId, time.Now().Unix())

	audioUrl, err := s.Data.TOS.Put(ctx, tos.PutRequest{
		Key:     fileName,
		Content: audioBytes,
	})
	if err != nil {
		return "", err
	}

	return audioUrl, nil
}
