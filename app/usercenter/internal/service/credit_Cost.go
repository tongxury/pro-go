package service

import (
	"context"
	"errors"
	creditpb "store/api/credit"
	"store/pkg/clients/mgz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *CreditService) XCost(ctx context.Context, request *creditpb.XCostRequest) (*creditpb.CreditState, error) {

	userId := request.UserId
	cost := request.Amount
	key := request.Key
	category := request.Category

	logger := log.NewHelper(log.With(log.DefaultLogger,
		"func", "XCost",
		"userId", userId,
		"cost", cost,
		"key", key,
		"category", category,
	))

	if cost == 0 {
		return nil, nil
	}

	// 找到未过期的所有记录
	recharges, err := t.data.Mongo.CreditRecharge.List(ctx,
		bson.M{
			"userId":   userId,
			"expireAt": bson.M{"$gt": time.Now().Unix()},
		},
		mgz.Find().Limit(1).B(),
	)
	if err != nil {
		return nil, err
	}

	if len(recharges) == 0 {
		return nil, errors.New("exceeded")
	}

	logger.Debugw("recharges", recharges)

	x := recharges[0]

	if x.Amount-x.Cost < request.Amount {
		return nil, errors.New("exceeded")
	}

	// 消耗
	_, err = t.data.Mongo.CreditRecharge.UpdateByIDIfExists(ctx,
		x.XId,
		mgz.Op().Inc("cost", request.Amount),
	)
	if err != nil {
		logger.Errorw("UpdateByIDIfExists err", err)
		return nil, err
	}

	_, _, err = t.data.Mongo.CreditChange.InsertNX(ctx,
		&creditpb.CreditChange{
			Amount:     request.Amount,
			UserId:     userId,
			RechargeId: x.XId,
			CreatedAt:  time.Now().Unix(),
			Uk:         key,
			Category:   category,
		},
		bson.M{"uk": key},
	)
	if err != nil {
		logger.Errorw("CreditChange.InsertNX err", err, "userId", userId)
	}

	return t.XGetCreditState(ctx, &creditpb.XGetCreditStateRequest{
		UserId: userId,
	})
}
