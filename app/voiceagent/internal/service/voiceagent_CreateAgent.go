package service

import (
	"context"
	ucpb "store/api/usercenter"
	voiceagent "store/api/voiceagent"
	"store/app/voiceagent/internal/data/repo/mongodb"
	"store/pkg/krathelper"
	"store/pkg/sdk/helper"
	"store/pkg/sdk/third/elevenlabs"
	"time"
)

func (s *VoiceAgentService) CreateAgent(ctx context.Context, req *voiceagent.CreateAgentRequest) (*voiceagent.Agent, error) {
	userId := krathelper.RequireUserId(ctx)

	var systemPrompt string
	voiceId := req.VoiceId
	name := req.Name
	avatar := req.Avatar
	desc := req.Desc

	var firstMessage string
	// 必须提供 personaId，从模板继承属性
	persona, err := s.Data.Mongo.Persona.GetById(ctx, req.PersonaId)
	if err != nil {
		return nil, err
	}
	if persona != nil {
		systemPrompt = mongodb.GenerateSystemPromptFromPersona(persona)
		if voiceId == "" {
			voiceId = persona.VoiceId
		}
		if name == "" {
			name = persona.DisplayName
		}
		if avatar == "" {
			avatar = persona.Avatar
		}
		if desc == "" {
			desc = persona.Description
		}

		firstMessage = persona.WelcomeMessage
	}

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

	elAgentId, err := s.Data.ElevenLabs.CreateAgent(ctx, agentReq)
	if err != nil {
		return nil, err
	}

	p := &voiceagent.Agent{
		User: &ucpb.User{
			XId: userId,
		},
		Persona: &voiceagent.Persona{
			XId:         req.PersonaId,
			DisplayName: name,
			Avatar:      avatar,
			Description: desc,
		},
		VoiceId:        voiceId,
		DefaultSceneId: req.DefaultSceneId,
		IsPublic:       req.IsPublic,
		Status:         "active",
		CreatedAt:      time.Now().Unix(),
		AgentId:        elAgentId,
	}

	res, err := s.Data.Mongo.Agent.Insert(ctx, p)
	if err != nil {
		_ = s.Data.ElevenLabs.DeleteAgent(ctx, elAgentId)
		return nil, err
	}
	return res, nil
}
