package mongodb

import (
	voiceagent "store/api/voiceagent"
	"store/pkg/clients/mgz"
	"go.mongodb.org/mongo-driver/mongo"
)

type PersonaCollection struct {
	*mgz.Core[*voiceagent.Persona]
}

func NewPersonaCollection(db *mongo.Database) *PersonaCollection {
	return &PersonaCollection{
		Core: mgz.NewCore[*voiceagent.Persona](db, "va_personas"),
	}
}
