package mongodb

import (
	"store/pkg/clients/mgz"
)

type Collections struct {
	Persona        *PersonaCollection
	Agent          *AgentCollection
	Voice          *VoiceCollection
	Scene          *SceneCollection
	Conversation   *ConversationCollection
	Transcript     *TranscriptCollection
	Motivation     *MotivationCollection
	Memory         *MemoryCollection
	UserProfile    *UserProfileCollection
	EmotionLog     *EmotionLogCollection
	ImportantEvent *ImportantEventCollection
	Assessment     *AssessmentCollection
}

func NewCollections(config mgz.Config) *Collections {
	database, err := mgz.Database(config)
	if err != nil {
		panic(err)
	}

	return &Collections{
		Persona:        NewPersonaCollection(database),
		Agent:          NewAgentCollection(database),
		Voice:          NewVoiceCollection(database),
		Scene:          NewSceneCollection(database),
		Conversation:   NewConversationCollection(database),
		Transcript:     NewTranscriptCollection(database),
		Motivation:     NewMotivationCollection(database),
		Memory:         NewMemoryCollection(database),
		UserProfile:    NewUserProfileCollection(database),
		EmotionLog:     NewEmotionLogCollection(database),
		ImportantEvent: NewImportantEventCollection(database),
		Assessment:     NewAssessmentCollection(database),
	}
}
