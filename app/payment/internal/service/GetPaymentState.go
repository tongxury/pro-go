package service

import (
	"context"
	paymentpb "store/api/payment"
	"store/app/payment/internal/data/repo/ent"
	"store/app/payment/internal/data/repo/ent/creditrecharge"
	"store/pkg/krathelper"
	"strings"
	"time"
)

func (t PaymentService) GetPaymentState(ctx context.Context, params *paymentpb.GetPaymentStateParams) (*paymentpb.PaymentState, error) {

	userId := krathelper.RequireUserId(ctx)

	recharges, err := t.data.Repos.EntClient.CreditRecharge.Query().
		Where(creditrecharge.UserID(userId)).
		Where(creditrecharge.StatusEQ("completed")).
		Where(creditrecharge.ExpireAtGT(time.Now())).
		Order(ent.Desc(creditrecharge.FieldID)).
		All(ctx)

	if err != nil {
		return nil, err
	}

	var subscribingPlanId string
	for _, x := range recharges {
		if strings.HasSuffix(x.PlanID, "ly") {
			subscribingPlanId = x.PlanID
			break
		}
	}

	return &paymentpb.PaymentState{
		SubscribingPlanId: subscribingPlanId,
	}, nil
}
