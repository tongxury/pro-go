package biz

import (
	"context"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/app/voiceagent/internal/data"
	"store/app/voiceagent/internal/data/repo/mongodb"
	"store/pkg/clients/mgz"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/elevenlabs"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

type AgentBiz struct {
	data *data.Data
}

func NewAgentBiz(data *data.Data) *AgentBiz {
	return &AgentBiz{
		data: data,
	}
}

func (b *AgentBiz) CreateDefaultCounselingAgent(ctx context.Context, userID string) error {
	// 使用内置的 "心心" (builtin_xin_xin) 角色
	name := "xinxin"
	persona, err := b.data.Mongo.Persona.FindOne(ctx, mgz.Filter().EQ("name", name).B())
	if err != nil {
		return err
	}
	if persona == nil {
		log.Infow("persona not found in mongo, getting from builtin memory", "name", name)
		personas := mongodb.GetBuiltinPersonas()
		for _, p := range personas {
			if p.Name == name {
				persona = p
				break
			}
		}
	}

	if persona == nil {
		log.Errorw("persona not found even in builtin memory", "name", name)
		return nil
	}

	systemPrompt := mongodb.GenerateSystemPromptFromPersona(persona)
	voiceId := persona.VoiceId
	firstMessage := persona.WelcomeMessage

	agentReq := &elevenlabs.CreateAgentRequest{
		Name: name,
		ConversationConfig: elevenlabs.ConversationConfig{
			Agent: elevenlabs.AgentConfig{
				Prompt: &elevenlabs.PromptSettings{
					Prompt: systemPrompt,
				},
				FirstMessage: helper.Pointer(firstMessage),
			},
			TTS: &elevenlabs.TTSConfig{
				VoiceID: voiceId,
			},
		},
	}

	elAgentId, err := b.data.ElevenLabs.CreateAgent(ctx, agentReq)
	if err != nil {
		return err
	}

	p := &voiceagent.Agent{
		User: &ucpb.User{
			XId: userID,
		},
		Persona: &voiceagent.Persona{
			XId:         persona.XId,
			Name:        persona.DisplayName,
			Avatar:      persona.Avatar,
			Description: persona.Description,
		},
		VoiceId:        voiceId,
		DefaultSceneId: "",
		IsPublic:       false,
		Status:         "active",
		CreatedAt:      time.Now().Unix(),
		AgentId:        elAgentId,
	}

	_, err = b.data.Mongo.Agent.Insert(ctx, p)
	if err != nil {
		_ = b.data.ElevenLabs.DeleteAgent(ctx, elAgentId)
		return err
	}

	return nil
}
