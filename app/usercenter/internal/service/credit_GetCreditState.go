package service

import (
	"context"
	creditpb "store/api/credit"
	"store/pkg/krathelper"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (t *CreditService) GetCreditState(ctx context.Context, request *creditpb.GetCreditStateRequest) (*creditpb.CreditState, error) {

	userId := krathelper.RequireUserId(ctx)

	//t.XCost(ctx, &creditpb.XCostRequest{
	//	UserId: "690d6a669e5c05462c0e4165",
	//	Amount: 10,
	//	Key:    "6933e58f1e6bf80a1b1e45e8",
	//})

	return t.XGetCreditState(ctx, &creditpb.XGetCreditStateRequest{
		UserId: userId,
	})
}

func (t *CreditService) XGetCreditState(ctx context.Context, request *creditpb.XGetCreditStateRequest) (*creditpb.CreditState, error) {

	userId := request.UserId

	recharges, err := t.data.Mongo.CreditRecharge.List(ctx,
		bson.M{
			"userId":   userId,
			"expireAt": bson.M{"$gt": time.Now().Unix()}},
		options.Find().SetSort(bson.M{"cost": -1}),
	)
	if err != nil {
		return nil, err
	}

	if len(recharges) == 0 {

		cat := creditpb.CreditRechargeCategoryNewUserReward
		const amount = 150

		newRecharge := &creditpb.CreditRecharge{
			Amount:    amount,
			UserId:    userId,
			CreatedAt: time.Now().Unix(),
			ExpireAt:  time.Now().Add(30 * 24 * time.Hour).Unix(),
			Category:  cat,
		}

		_, _, _ = t.data.Mongo.CreditRecharge.InsertNX(ctx, newRecharge,
			bson.M{"uk": cat + "_" + userId},
		)

		return &creditpb.CreditState{
			UserId:        userId,
			Total:         amount,
			Balance:       amount,
			LastUpdatedAt: time.Now().Unix(),
		}, nil
	}

	// 可用余额
	var total int64
	var balance int64
	for _, x := range recharges {
		total += x.Amount
		balance += x.Amount - x.Cost
	}

	return &creditpb.CreditState{
		UserId:  userId,
		Total:   total,
		Balance: balance,
	}, nil
}
