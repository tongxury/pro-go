package service

import (
	"store/api/voiceagent"
	"store/app/voiceagent/internal/data"
)

type LiveKitService struct {
	voiceagent.UnimplementedLiveKitServiceServer
	data *data.Data
}

func NewLiveKitService(data *data.Data) *LiveKitService {
	return &LiveKitService{
		data: data,
	}
}
