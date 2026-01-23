package service

import (
	"context"
	paymentpb "store/api/payment"
	"store/app/payment/internal/data/repo/ent/payment"
	"store/pkg/krathelper"
)

func (t PaymentService) ListPayments(ctx context.Context, params *paymentpb.ListPaymentsParams) (*paymentpb.PaymentList, error) {

	userId := krathelper.RequireUserId(ctx)

	if params.Page == 0 {
		params.Page = 1
	}

	t.data.Repos.EntClient.Payment.Query().
		Where(payment.UserID(userId)).
		Limit(int(params.Size)).
		Offset(int((params.Page - 1) * params.Size)).
		All(ctx)

	return &paymentpb.PaymentList{}, nil
}
