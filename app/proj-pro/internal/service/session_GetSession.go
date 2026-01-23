package service

import (
	"context"
	projpb "store/api/proj"

	"go.mongodb.org/mongo-driver/bson"
)

func (t *SessionService) GetSession(ctx context.Context, params *projpb.GetSessionRequest) (*projpb.Session, error) {
	task, err := t.data.Mongo.Session.GetById(ctx, params.Id)
	if err != nil {
		return nil, err
	}

	segments, err := t.data.Mongo.SessionSegment.List(ctx, bson.M{"sessionId": task.XId})
	if err != nil {
		return nil, err
	}

	task.Segments = segments

	return task, nil
}
