package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type SessionSegmentCollection struct {
	*mgz.Core[*projpb.SessionSegment]
}

func NewSessionSegmentCollection(db *mongo.Database) *SessionSegmentCollection {
	return &SessionSegmentCollection{
		Core: mgz.NewCore[*projpb.SessionSegment](db, "session_segments"),
	}
}
