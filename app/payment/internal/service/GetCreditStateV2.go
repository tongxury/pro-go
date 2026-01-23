package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	paymentpb "store/api/payment"
	"store/pkg/krathelper"
)

func (t PaymentService) GetCreditStateV2(ctx context.Context, params *paymentpb.GetCreditStateV2Params) (*paymentpb.CreditState, error) {

	userId := krathelper.RequireUserId(ctx)

	state, err := t.credit.GetCreditState(ctx, userId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return state, nil
}
