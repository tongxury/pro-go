package repo

import (
	projpb "store/api/proj"
	"store/pkg/clients/mgz"

	"go.mongodb.org/mongo-driver/mongo"
)

type FeedbackCollection struct {
	*mgz.Core[*projpb.Feedback]
}

func NewFeedbackCollection(db *mongo.Database) *FeedbackCollection {
	return &FeedbackCollection{
		Core: mgz.NewCore[*projpb.Feedback](db, "feedbacks"),
	}
}
