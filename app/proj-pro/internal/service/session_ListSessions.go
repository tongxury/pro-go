package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/krathelper"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (t *SessionService) ListSessions(ctx context.Context, req *projpb.ListSessionsRequest) (*projpb.SessionList, error) {
	userId := krathelper.RequireUserId(ctx)

	list, _, err := t.data.Mongo.Session.ListAndCount(ctx,
		bson.M{
			"userId": userId,
		},
		options.Find().SetSort(bson.M{"createdAt": -1}),
	)
	if err != nil {
		return nil, err
	}

	return &projpb.SessionList{
		List: list,
	}, nil
}
