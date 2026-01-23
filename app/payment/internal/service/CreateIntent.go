package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	paymentpb "store/api/payment"
	"store/pkg/krathelper"
)

func (t PaymentService) CreateIntent(ctx context.Context, params *paymentpb.CreateIntentParams) (*paymentpb.Intent, error) {

	userId := krathelper.RequireUserId(ctx)

	intent, err := t.payment.CreateIntent(ctx, userId, params)
	if err != nil {
		log.Errorw("create intent error", err, "userId", userId)
		return nil, err
	}

	return intent, nil
}
