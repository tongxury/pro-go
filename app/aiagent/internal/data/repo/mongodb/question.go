package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	aiagentpb "store/api/aiagent"
	"store/pkg/clients/mgz"
)

type QuestionCollection struct {
	*mgz.CoreCollection[aiagentpb.Question]
}

func NewQuestionCollection(db *mongo.Database) *QuestionCollection {
	return &QuestionCollection{
		CoreCollection: mgz.NewCoreCollection[aiagentpb.Question](db, "questions"),
	}
}

func (t *QuestionCollection) InsertIfNX(ctx context.Context, filters bson.M, data *aiagentpb.Question) (*aiagentpb.Question, error) {

	olds, err := t.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	if len(olds) > 0 {
		return olds[0], nil
	}

	err = t.Insert(ctx, data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
