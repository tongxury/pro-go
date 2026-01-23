package mongodb

import (
	"store/pkg/clients/mgz"
)

type Collections struct {
	Persona      *PersonaCollection
	Agent        *AgentCollection
	Voice        *VoiceCollection
	Scene        *SceneCollection
	Conversation *ConversationCollection
	Transcript   *TranscriptCollection
	Motivation   *MotivationCollection
}

func NewCollections(config mgz.Config) *Collections {
	database, err := mgz.Database(config)
	if err != nil {
		panic(err)
	}

	return &Collections{
		Persona:      NewPersonaCollection(database),
		Agent:        NewAgentCollection(database),
		Voice:        NewVoiceCollection(database),
		Scene:        NewSceneCollection(database),
		Conversation: NewConversationCollection(database),
		Transcript:   NewTranscriptCollection(database),
		Motivation:   NewMotivationCollection(database),
	}
}
