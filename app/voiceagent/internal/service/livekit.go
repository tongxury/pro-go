package service

import (
	"store/api/voiceagent"
)

type LiveKitService struct {
	voiceagent.UnimplementedLiveKitServiceServer
}

func NewLiveKitService() *LiveKitService {
	return &LiveKitService{}
}
