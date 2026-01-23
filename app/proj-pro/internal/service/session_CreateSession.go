package service

import (
	"context"
	projpb "store/api/proj"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func (t *SessionService) CreateSession(ctx context.Context, params *projpb.CreateSessionRequest) (*projpb.Session, error) {

	userId := krathelper.RequireUserId(ctx)

	comm, err := t.data.Mongo.Commodity.FindByID(ctx, params.CommodityId)
	if err != nil {
		return nil, err
	}

	targetChanceIndex := 0

	newSession := &projpb.Session{
		Commodity:    comm,
		TargetChance: comm.Chances[targetChanceIndex],
		Status:       "chanceSelected",
		CreatedAt:    time.Now().Unix(),
		UserId:       userId,
	}

	session, _, err := t.data.Mongo.Session.InsertNX(ctx,
		newSession,
		bson.M{"userId": userId,
			"commodity._id":      params.CommodityId,
			"targetChance.index": targetChanceIndex,
		},
	)
	if err != nil {
		return nil, err
	}

	newSession.XId = session

	return newSession, nil
}
