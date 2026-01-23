package service

import (
	"context"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	paymentpb "store/api/payment"
	"store/app/payment/configs"
	"store/app/payment/internal/data/repo/ent/creditrecharge"
	"store/app/payment/internal/data/repo/ent/payment"
	"store/pkg/enums"
	"store/pkg/events"
	"store/pkg/sdk/conv"
	"time"
)

func (t PaymentService) CallbackByAirWallex(ctx context.Context, params *paymentpb.CallbackByAirWallexParams) (*empty.Empty, error) {

	log.Debugw("CallbackByAirWallex ", "", "params", params)

	if params.Name == "payment_intent.succeeded" {

		paymentId := params.GetData().GetObject().GetMerchantOrderId()

		if paymentId == "" {
			return nil, errors.BadRequest("empty paymentId", "")
		}

		err := t.data.Repos.EntClient.Payment.Update().
			SetStatus(enums.PaymentStatus_Complete).
			SetExpireAt(time.Now().Add(30 * 24 * time.Hour)).
			Where(payment.ID(conv.Int64(paymentId))).
			Exec(ctx)

		if err != nil {
			return nil, err
		}

		// 双写
		err = t.data.Repos.EntClient.CreditRecharge.Update().
			SetStatus("completed").
			SetExpireAt(time.Now().Add(30 * 24 * time.Hour)).
			Where(creditrecharge.Key(paymentId)).
			Exec(ctx)

		if err != nil {
			log.Errorw("CreditRecharge update err", err)
		}

		t.data.Repos.LocalCache.Flush()

		pay, err := t.data.Repos.EntClient.Payment.Get(ctx, conv.Int64(paymentId))
		if err != nil {
			return nil, err
		}

		plan := configs.GetPlanById(pay.PlanID)

		t.productEvent(ctx, events.PaymentSuccessEvent{
			UniqueId: paymentId,
			UserID:   pay.UserID,
			//DeviceID: deviceId,
			Amount:   plan.Amount,
			Platform: "airwallex",
			Ts:       time.Now().Unix(),
		})
	}

	return &empty.Empty{}, nil
}
