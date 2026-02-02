package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"store/pkg/sdk/third/elevenlabs"
)

func (s *VoiceAgentService) UpdateAgent(ctx context.Context, req *voiceagent.UpdateAgentRequest) (*voiceagent.Agent, error) {
	p, err := s.Data.Mongo.Agent.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// 如果模板 ID 发生变化，更新系统提示词
	if req.PersonaId != "" && (p.Persona == nil || req.PersonaId != p.Persona.XId) {
		persona, err := s.Data.Mongo.Persona.GetById(ctx, req.PersonaId)
		if err == nil && persona != nil {
			p.Persona = persona
		}
	}

	if p.AgentId != "" {
		updateBody := map[string]any{
			"name": req.Name,
			"conversation_config": elevenlabs.ConversationConfig{
				Agent: elevenlabs.AgentConfig{
					Prompt: &elevenlabs.PromptSettings{
						//Text: p.SystemPrompt,
					},
				},
				TTS: &elevenlabs.TTSConfig{
					VoiceID: req.VoiceId,
				},
			},
		}
		if err := s.Data.ElevenLabs.UpdateAgent(ctx, p.AgentId, updateBody); err != nil {
			return nil, err
		}
	}

	if p.Persona == nil {
		p.Persona = &voiceagent.Persona{}
	}
	p.Persona.DisplayName = req.Name
	p.Persona.Avatar = req.Avatar
	p.Persona.Description = req.Desc
	p.VoiceId = req.VoiceId
	p.DefaultSceneId = req.DefaultSceneId
	p.IsPublic = req.IsPublic
	p.Status = req.Status

	if _, err := s.Data.Mongo.Agent.ReplaceByID(ctx, p.XId, p); err != nil {
		return nil, err
	}

	return p, nil
}
