package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type AgentCollection struct {
	*mgz.Core[*voiceagent.Agent]
}

func NewAgentCollection(db *mongo.Database) *AgentCollection {
	return &AgentCollection{
		Core: mgz.NewCore[*voiceagent.Agent](db, "va_agents"),
	}
}
