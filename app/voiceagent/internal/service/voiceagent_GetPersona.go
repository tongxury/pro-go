package service

import (
	"context"
	voiceagent "store/api/voiceagent"
)

func (s *VoiceAgentService) GetPersona(ctx context.Context, req *voiceagent.GetPersonaRequest) (*voiceagent.Persona, error) {
	return s.Data.Mongo.Persona.GetById(ctx, req.Id)
}
