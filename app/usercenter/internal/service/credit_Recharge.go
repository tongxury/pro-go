package service

import (
	"context"
	creditpb "store/api/credit"
	"store/pkg/clients/mgz"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"go.mongodb.org/mongo-driver/bson"
)

func (t *CreditService) XRecharge(ctx context.Context, request *creditpb.XRechargeRequest) (*creditpb.CreditState, error) {

	userId := request.UserId
	amount := request.Amount
	key := request.Key
	category := request.Category

	if amount == 0 {
		return nil, nil
	}

	// 找到未过期的所有记录 正常只有1个
	var inheritedAmount int64
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

	if len(recharges) > 0 {
		left := recharges[0].Amount - recharges[0].Cost

		if left > 0 {
			inheritedAmount = left
		}
	}

	id, _, err := t.data.Mongo.CreditRecharge.InsertNX(ctx,
		&creditpb.CreditRecharge{
			Amount:          amount + inheritedAmount,
			InheritedAmount: inheritedAmount,
			UserId:          userId,
			CreatedAt:       time.Now().Unix(),
			ExpireAt:        time.Now().Add(30 * 24 * time.Hour).Unix(),
			Uk:              key,
			Category:        category,
		},
		mgz.Filter().EQ("key", key).B(),
	)

	if err != nil {
		log.Errorw("InsertNX err", err)
		return nil, err
	}

	_, _, err = t.data.Mongo.CreditChange.InsertNX(ctx,
		&creditpb.CreditChange{
			Amount:     request.Amount,
			UserId:     userId,
			RechargeId: id,
			CreatedAt:  time.Now().Unix(),
			Uk:         key,
			Category:   category,
		},
		bson.M{"uk": key},
	)
	if err != nil {
		log.Errorw("CreditChange.InsertNX err", err, "userId", userId)
	}

	return t.XGetCreditState(ctx, &creditpb.XGetCreditStateRequest{
		UserId: userId,
	})
}
