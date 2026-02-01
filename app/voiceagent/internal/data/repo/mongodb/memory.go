package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type MemoryCollection struct {
	*mgz.Core[*voiceagent.Memory]
}

func NewMemoryCollection(db *mongo.Database) *MemoryCollection {
	return &MemoryCollection{
		Core: mgz.NewCore[*voiceagent.Memory](db, "va_memories"),
	}
}
