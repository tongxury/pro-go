package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/app/voiceagent/internal/biz"
	"store/app/voiceagent/internal/data"
)

type VoiceAgentService struct {
	voiceagent.UnimplementedVoiceAgentServiceServer
	voiceagent.UnimplementedMotivationServiceServer
	Data     *data.Data
	item     *biz.ItemBiz
	AgentBiz *biz.AgentBiz
	lk       *LiveKitService
}

func NewVoiceAgentService(data *data.Data, item *biz.ItemBiz, agentBiz *biz.AgentBiz, lk *LiveKitService) *VoiceAgentService {
	s := &VoiceAgentService{
		Data:     data,
		item:     item,
		AgentBiz: agentBiz,
		lk:       lk,
	}
	// 自动启动异步同步循环
	//s.StartSyncLoop(context.Background())

	s.SeedSystemVoices(context.Background())
	s.SyncPersonas(context.Background())
	return s
}
