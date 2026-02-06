package service

import (
	"store/api/voiceagent"
	"store/app/voiceagent/internal/biz"
	"store/app/voiceagent/internal/data"
)

type LiveKitService struct {
	voiceagent.UnimplementedLiveKitServiceServer
	data     *data.Data
	agentBiz *biz.AgentBiz
}

func NewLiveKitService(data *data.Data, agentBiz *biz.AgentBiz) *LiveKitService {
	return &LiveKitService{
		data:     data,
		agentBiz: agentBiz,
	}
}
