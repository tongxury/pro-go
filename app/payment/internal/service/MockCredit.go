package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	paymentpb "store/api/payment"
	userpb "store/api/user"
	"store/app/payment/configs"
	"store/app/payment/internal/data/repo/ent"
	"store/app/payment/internal/data/repo/ent/creditrecharge"
	"strings"
	"time"
)

func (t PaymentService) MockCredit(ctx context.Context, params *paymentpb.MockCreditParams) (*paymentpb.CreditState, error) {

	user, err := t.data.GrpcClients.UserClient.GetUserByEmailOrPhone(ctx, &userpb.GetUserByEmailOrPhoneParams{
		Value: params.EmailOrPhone,
	})

	if err != nil {
		return nil, err
	}

	switch params.Action {
	case "recharge":

		if strings.HasSuffix(params.PlanId, "ly") {
			return nil, errors.BadRequest("invalidPlan", "")
		}

		plan := configs.GetPlanById(params.PlanId)
		if plan == nil {
			return nil, errors.BadRequest("invalidPlan", "")
		}

		err = t.data.Repos.EntClient.CreditRecharge.Create().
			SetStatus("completed").
			SetPlatform("manual").
			SetPlanID(params.PlanId).
			SetQuota(plan.CreditPerMonth).
			SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
			SetUserID(user.Id).
			SetKey(fmt.Sprintf("%s_%d", user.Id, time.Now().Unix())).
			OnConflictColumns(creditrecharge.FieldKey).
			DoNothing().
			Exec(ctx)

		if err != nil {
			log.Error("CreditChange Create Error", "userId", user.Id, "planId", plan.Id)
			return nil, err
		}

	case "cancel":
		_, err = t.data.Repos.EntClient.CreditRecharge.Delete().
			Where(creditrecharge.UserID(user.Id)).
			Exec(ctx)

		if err != nil {
			log.Error("CreditChange Delete Error", "userId", user.Id, "planId")
			return nil, err
		}

	case "incr":

		all, err := t.data.Repos.EntClient.CreditRecharge.Query().
			Where(creditrecharge.UserID(user.Id)).
			Where(creditrecharge.Status("completed")).
			Where(creditrecharge.ExpireAtGT(time.Now())).
			Order(ent.Desc(creditrecharge.FieldID)).
			All(ctx)
		if err != nil {
			return nil, err
		}

		if len(all) == 0 {
			return nil, nil
		}

		err = t.data.Repos.EntClient.CreditRecharge.Update().Where(creditrecharge.ID(all[0].ID)).AddQuota(params.Amount).Exec(ctx)
		if err != nil {
			return nil, err
		}
	}

	return t.credit.GetCreditState(ctx, user.Id)
}

//func (t PaymentService) MockDecrCredit(ctx context.Context, params *paymentpb.MockDecrCreditParams) (*paymentpb.CreditState, error) {
//	user, err := t.data.GrpcClients.UserClient.GetUserByEmailOrPhone(ctx, &userpb.GetUserByEmailOrPhoneParams{
//		Value: params.EmailOrPhone,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	ongoingPayment, err := t.payment.GetOngoingPayment(ctx, user.Id)
//	if err != nil {
//		return nil, err
//	}
//
//	if ongoingPayment == nil {
//		return nil, nil
//	}
//
//	err = t.data.Repos.RedisClient.DecrBy(ctx, fmt.Sprintf("credit:used.%d_%s", ongoingPayment.ID, user.Id), params.Amount).Err()
//	if err != nil {
//		return nil, err
//	}
//
//	state, err := t.payment.GetCreditState(ctx, user.Id)
//	if err != nil {
//		return nil, err
//	}
//
//	return state, nil
//}
//
//func (t PaymentService) MockDiscardPayment(ctx context.Context, params *paymentpb.MockDiscardPaymentParams) (*paymentpb.CreditState, error) {
//
//	user, err := t.data.GrpcClients.UserClient.GetUserByEmailOrPhone(ctx, &userpb.GetUserByEmailOrPhoneParams{
//		Value: params.EmailOrPhone,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	err = t.data.Repos.EntClient.Payment.Update().
//		Where(payment.UserID(user.Id)).
//		SetExpireAt(time.Now()).
//		Exec(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	t.data.Repos.LocalCache.Delete("ongoingPayment:" + user.Id)
//
//	state, err := t.payment.GetCreditState(ctx, user.Id)
//	if err != nil {
//		return nil, err
//	}
//
//	return state, nil
//}
//
//func (t PaymentService) MockPayment(ctx context.Context, params *paymentpb.MockPaymentParams) (*paymentpb.MockPaymentResult, error) {
//
//	user, err := t.data.GrpcClients.UserClient.GetUserByEmailOrPhone(ctx, &userpb.GetUserByEmailOrPhoneParams{
//		Value: params.EmailOrPhone,
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	plan := configs.GetPlanById(params.PlanId)
//
//	if plan == nil {
//		return nil, errors.BadRequest("invalidPlan", "")
//	}
//
//	pay, err := t.data.Repos.EntClient.Payment.Create().
//		SetStatus(enums.PaymentStatus_Complete).
//		SetPlatform("manual").
//		SetAmount(plan.CnyAmount).
//		SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
//		SetPlanID(params.PlanId).
//		SetUserID(user.Id).
//		//OnConflictColumns(payment.FieldPlatform, payment.FieldSubID).
//		//UpdateExpireAt().
//		Save(ctx)
//	if err != nil {
//		return nil, err
//	}
//
//	t.productEvent(ctx, events.PaymentSuccessEvent{
//		UniqueId: conv.Str(pay.ID),
//		UserID:   user.Id,
//		//DeviceID: deviceId,
//		Amount:   plan.Amount,
//		Platform: "manual",
//	})
//
//	t.data.Repos.LocalCache.Delete("ongoingPayment:" + user.Id)
//
//	state, err := t.payment.GetCreditState(ctx, user.Id)
//	if err != nil {
//		return nil, err
//	}
//
//	return &paymentpb.MockPaymentResult{
//		PlanId:    params.PlanId,
//		Remaining: state.Total - state.Used,
//		PaymentId: conv.Str(pay.ID),
//		UserId:    user.Id,
//	}, nil
//
//}
