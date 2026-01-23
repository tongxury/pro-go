package service

import (
	voiceagent "store/api/voiceagent"
	"store/app/voiceagent/internal/biz"
	"store/app/voiceagent/internal/data"
)

type VoiceAgentService struct {
	voiceagent.UnimplementedVoiceAgentServiceServer
	voiceagent.UnimplementedMotivationServiceServer
	Data *data.Data
	item *biz.ItemBiz
}

func NewVoiceAgentService(data *data.Data, item *biz.ItemBiz) *VoiceAgentService {
	s := &VoiceAgentService{
		Data: data,
		item: item,
	}
	// 自动启动异步同步循环
	//s.StartSyncLoop(context.Background())
	return s
}
