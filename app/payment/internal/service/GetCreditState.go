package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
	paymentpb "store/api/payment"
)

func (t PaymentService) GetCreditState(ctx context.Context, params *paymentpb.GetCreditStateParams) (*paymentpb.CreditState, error) {

	state, err := t.credit.GetCreditState(ctx, params.UserId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return state, nil
}
