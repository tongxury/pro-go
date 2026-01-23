package service

import (
	"context"
	voiceagent "store/api/voiceagent"
	"go.mongodb.org/mongo-driver/bson"
)

func (s *VoiceAgentService) ListPersonas(ctx context.Context, req *voiceagent.ListPersonasRequest) (*voiceagent.PersonaList, error) {
	filter := bson.M{}
	if req.Category != "" {
		filter["category"] = req.Category
	}

	list, err := s.Data.Mongo.Persona.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return &voiceagent.PersonaList{
		List: list,
	}, nil
}
