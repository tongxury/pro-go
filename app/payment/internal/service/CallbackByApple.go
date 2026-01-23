package service

import (
	"context"
	"github.com/awa/go-iap/appstore"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang/protobuf/ptypes/empty"
	paymentpb "store/api/payment"
	"store/app/payment/configs"
	"store/app/payment/internal/data/repo/ent/creditrecharge"
	"store/app/payment/internal/data/repo/ent/payment"
	"store/pkg/enums"
	"store/pkg/events"
	"store/pkg/krathelper"
	"store/pkg/sdk/conv"
	"strings"
	"time"
)

func (t PaymentService) CallbackByApple(ctx context.Context, params *paymentpb.CallbackByAppleParams) (*empty.Empty, error) {

	deviceId := krathelper.GetDeviceId(ctx)

	userId := krathelper.RequireUserId(ctx)

	var receipt struct {
		TransactionId string
		ProductId     string
	}

	_ = conv.J2S([]byte(params.TransactionReceipt), &receipt)

	log.Debugw("callback  by apple", params, "userId", userId, "receipt", receipt)

	client := appstore.New()

	req := appstore.IAPRequest{
		ReceiptData: params.TransactionReceipt,
	}

	var resp appstore.IAPResponse

	err := client.Verify(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	log.Debugw("callback resp by apple", conv.S2J(resp))

	//planId := strings.ReplaceAll(params.ProductId, "com_veogo_", "")
	planId := receipt.ProductId
	//planId = strings.ReplaceAll(planId, "_", "-")

	plan := configs.GetPlanById(planId)

	if plan == nil {
		log.Warnw("callback by apple", "planId", planId)
		if strings.Contains(planId, "l1") {
			plan = configs.GetPlanById("l1-month")
		}
		if strings.Contains(planId, "l2") {
			plan = configs.GetPlanById("l2-month")
		}
		if strings.Contains(planId, "l3") {
			plan = configs.GetPlanById("l3-month")
		}

		log.Warnw("callback by apple1 ", "planId", planId, "plan", plan)
	}

	if plan == nil {
		return nil, errors.BadRequest("invalidPlan:"+planId, "")
	}

	planId = plan.Id

	err = t.data.Repos.EntClient.Payment.Create().
		SetStatus(enums.PaymentStatus_Complete).
		SetPlatform("apple").
		SetAmount(plan.Amount).
		SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
		SetPlanID(planId).
		SetSessionID(receipt.TransactionId).
		SetUserID(userId).
		OnConflictColumns(payment.FieldSessionID).
		DoNothing().
		//UpdateExpireAt().
		Exec(ctx)

	if err != nil {
		log.Error("Payment Create Error", err, "userId", userId, "planId", plan.Id)
		return nil, err
	}

	// 双写
	err = t.data.Repos.EntClient.CreditRecharge.Create().
		SetStatus("completed").
		SetPlatform("apple").
		SetPlanID(planId).
		SetQuota(plan.CreditPerMonth).
		SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
		SetUserID(userId).
		SetKey(receipt.TransactionId).
		OnConflictColumns(creditrecharge.FieldKey).
		DoNothing().
		Exec(ctx)
	if err != nil {
		log.Error("CreditChange Create Error", "userId", userId, "planId", plan.Id)
	}

	t.data.Repos.LocalCache.Delete("ongoingPayment:" + userId)

	t.productEvent(ctx, events.PaymentSuccessEvent{
		UniqueId: receipt.TransactionId,
		UserID:   userId,
		DeviceID: deviceId,
		Amount:   plan.Amount,
		Platform: "apple",
		Ts:       time.Now().Unix(),
	})

	return &empty.Empty{}, nil
}
