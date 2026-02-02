package service

import (
	"context"
	"fmt"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/pkg/sdk/third/elevenlabs"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// StartSyncLoop 启动定时同步任务
func (s *VoiceAgentService) StartSyncLoop(ctx context.Context) {

	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		// 启动后立即执行一次
		s.SyncData(ctx)
		for {
			select {
			case <-ticker.C:
				s.SyncData(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}

// SyncData 执行全量数据同步 (Agent & Voice & Persona)
func (s *VoiceAgentService) SyncData(ctx context.Context) {
	fmt.Println("[Sync] Starting data synchronization with ElevenLabs...")
	s.SyncAgents(ctx)
	//s.SyncVoices(ctx)
	fmt.Println("[Sync] Data synchronization completed.")
}

// SyncPersonas 同步并初始化系统内置角色模板
func (s *VoiceAgentService) SyncPersonas(ctx context.Context) {
	err := s.Data.Mongo.Persona.SeedBuiltinPersonas(ctx)
	if err != nil {
		fmt.Printf("[Sync] Failed to seed builtin personas: %v\n", err)
	} else {
		fmt.Println("[Sync] Builtin personas synchronized.")
	}
}

// SyncAgents 同步 Agent 数据
func (s *VoiceAgentService) SyncAgents(ctx context.Context) {
	// 1. 获取远程所有 Agents
	resp, err := s.Data.ElevenLabs.ListAgents(ctx, &elevenlabs.ListAgentsParams{PageSize: 100})
	if err != nil {
		fmt.Printf("[Sync] Failed to list remote agents: %v\n", err)
		return
	}

	remoteAgents := make(map[string]elevenlabs.Agent)
	for _, a := range resp.Agents {
		remoteAgents[a.AgentID] = a
	}

	// 2. 获取本地所有 Agents
	localAgents, err := s.Data.Mongo.Agent.List(ctx, bson.M{})
	if err != nil {
		fmt.Printf("[Sync] Failed to list local agents: %v\n", err)
		return
	}

	// 3. 补充或更新
	for id, ra := range remoteAgents {
		found := false
		for _, la := range localAgents {
			if la.AgentId == id {
				// 更新逻辑
				if la.Persona == nil {
					la.Persona = &voiceagent.Persona{}
				}
				la.Persona.DisplayName = ra.Name

				if ra.ConversationConfig.Agent.Prompt != nil {
					//la.SystemPrompt = ra.ConversationConfig.Agent.Prompt.Text
				}

				_, _ = s.Data.Mongo.Agent.ReplaceByID(ctx, la.XId, la)
				found = true
				break
			}
		}

		if !found {
			// 本地缺失，补充

			//var prompt string
			//if ra.ConversationConfig.Agent.Prompt != nil {
			//	prompt = ra.ConversationConfig.Agent.Prompt.Text
			//}

			newAgent := &voiceagent.Agent{
				//XId:          primitive.NewObjectID().Hex(),
				User: &ucpb.User{
					XId: "system", // 同步的数据标记为系统
				},
				Persona: &voiceagent.Persona{
					DisplayName: ra.Name,
				},
				//SystemPrompt: prompt,
				AgentId:   ra.AgentID,
				Status:    "active",
				CreatedAt: ra.Metadata.CreatedAtUnixSecs,
			}
			_, _ = s.Data.Mongo.Agent.Insert(ctx, newAgent)
			fmt.Printf("[Sync] Added missing agent: %s\n", ra.Name)
		}
	}

	// 4. 删除本地多余 (在该 SDK 中已不存在的)
	for _, la := range localAgents {
		if la.AgentId != "" {
			if _, exists := remoteAgents[la.AgentId]; !exists {
				// 远程已删除，本地同步物理删除或软删除 (此处根据要求执行物理删除)
				_ = s.Data.Mongo.Agent.DeleteByID(ctx, la.XId)
				name := ""
				if la.Persona != nil {
					name = la.Persona.DisplayName
				}
				fmt.Printf("[Sync] Deleted redundant agent: %s\n", name)
			}
		}
	}
}

// SyncVoices 同步 Voice 数据
func (s *VoiceAgentService) SyncVoices(ctx context.Context) {
	// 1. 获取远程所有 Voices
	resp, err := s.Data.ElevenLabs.ListVoices(ctx)
	if err != nil {
		fmt.Printf("[Sync] Failed to list remote voices: %v\n", err)
		return
	}

	remoteVoices := make(map[string]elevenlabs.Voice)
	for _, v := range resp.Voices {
		remoteVoices[v.VoiceID] = v
	}

	// 2. 获取本地所有 Voices
	localVoices, err := s.Data.Mongo.Voice.List(ctx, bson.M{})
	if err != nil {
		fmt.Printf("[Sync] Failed to list local voices: %v\n", err)
		return
	}

	// 3. 补充或更新
	for id, rv := range remoteVoices {
		found := false
		for _, lv := range localVoices {
			if lv.VoiceId == id {
				lv.Name = rv.Name
				lv.Type = rv.Category
				_, _ = s.Data.Mongo.Voice.ReplaceByID(ctx, lv.XId, lv)
				found = true
				break
			}
		}

		if !found {
			newVoice := &voiceagent.Voice{
				XId:       primitive.NewObjectID().Hex(),
				UserId:    "system",
				Name:      rv.Name,
				Type:      rv.Category,
				VoiceId:   rv.VoiceID,
				Status:    "active",
				CreatedAt: time.Now().Unix(),
			}
			_, _ = s.Data.Mongo.Voice.Insert(ctx, newVoice)
			fmt.Printf("[Sync] Added missing voice: %s\n", rv.Name)
		}
	}

	// 4. 删除本地多余
	for _, lv := range localVoices {
		if lv.VoiceId != "" {
			if _, exists := remoteVoices[lv.VoiceId]; !exists {
				_ = s.Data.Mongo.Voice.DeleteByID(ctx, lv.XId)
				fmt.Printf("[Sync] Deleted redundant voice: %s\n", lv.Name)
			}
		}
	}
}
