package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type SessionCollection struct {
	*mgz.Core[*projpb.Session]
}

func NewSessionCollection(db *mongo.Database) *SessionCollection {
	return &SessionCollection{
		Core: mgz.NewCore[*projpb.Session](db, "sessions"),
	}
}
