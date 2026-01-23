package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	paymentpb "store/api/payment"
	userpb "store/api/user"
	"store/app/payment/configs"
	"store/pkg/enums"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/third/airwallex"
)

func (t *PaymentBiz) CreateIntent(ctx context.Context, userId string, params *paymentpb.CreateIntentParams) (*paymentpb.Intent, error) {

	user, err := t.data.GrpcClients.UserClient.GetUserById(ctx, &userpb.GetUserByIdParams{Id: userId})
	if err != nil {
		return nil, err
	}

	planId := params.GetPlanId()
	if planId == "" {
		planId = fmt.Sprintf("%s-%s", params.Level, params.Cycle)
	}

	plan := configs.GetPlanById(planId)

	pay, err := t.data.Repos.EntClient.Payment.Create().
		SetStatus(enums.PaymentStatus_Created).
		SetPlatform("airwallex").
		SetAmount(plan.Amount).
		SetPlanID(planId).
		SetUserID(user.Id).
		//OnConflictColumns(payment.FieldPlatform, payment.FieldSubID).
		//UpdateExpireAt().
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// 双写
	_, err = t.data.Repos.EntClient.CreditRecharge.Create().
		SetStatus("pending").
		SetPlatform("airwallex").
		SetPlanID(planId).
		SetQuota(plan.CreditPerMonth).
		SetUserID(user.Id).
		SetKey(conv.Str(pay.ID)).
		Save(ctx)
	if err != nil {
		log.Errorw("CreditChange Create err", err, "id", pay.ID)
	}

	customer, err := t.data.AirWallex.GetOrCreateCustomer(ctx, airwallex.CreateCustomerParams{
		MerchantCustomerId: user.Id,
		FirstName:          user.Username,
		//LastName:           "",
		Email:       user.Email,
		PhoneNumber: user.Phone,
		Metadata:    nil,
	})
	if err != nil {
		return nil, err
	}

	log.Debugw("customer", customer)

	intent, err := t.data.AirWallex.CreatePaymentIntent(ctx, airwallex.CreatePaymentIntentParams{
		Amount: 0,
		//Amount:          plan.Amount,
		MerchantOrderId: conv.Str(pay.ID),
		Product: airwallex.CreatePaymentIntentParams_Product{
			Name: plan.Title,
			//Description: "test test",
			Quantity:  1,
			UnitPrice: plan.Amount,
		},
		CustomerId: customer.Id,
		Currency:   "CNY",
		ReturnUrl:  params.SuccessUrl,
		Metadata:   map[string]any{},
	})
	if err != nil {
		return nil, err
	}

	return &paymentpb.Intent{
		ClientSecret: intent.ClientSecret,
		Amount:       float64(intent.Amount),
		Id:           intent.Id,
		Currency:     intent.Currency,
		Mode:         plan.Mode,
		CountryCode:  "CN",
	}, nil
}
