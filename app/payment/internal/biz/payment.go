package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
	"github.com/stripe/stripe-go/v82"
	paymentpb "store/api/payment"
	userpb "store/api/user"
	"store/app/payment/configs"
	"store/app/payment/internal/data"
	"store/app/payment/internal/data/repo/ent"
	"store/app/payment/internal/data/repo/ent/payment"
	"store/app/payment/internal/data/stripefactory"
	"store/pkg/enums"
	"store/pkg/sdk/conv"
	"store/pkg/sdk/helper"
	"store/pkg/types"
	"strings"
	"time"
)

type PaymentBiz struct {
	data *data.Data
}

func NewPaymentBiz(data *data.Data) *PaymentBiz {
	return &PaymentBiz{data: data}
}

func (t *PaymentBiz) GetCreditState(ctx context.Context, userId string) (*paymentpb.CreditState, error) {

	ongoingPayment, err := t.GetOngoingPayment(ctx, userId)
	if err != nil {
		return nil, err
	}

	if ongoingPayment == nil {

		creditUsed, err := t.data.Repos.RedisClient.Get(ctx, fmt.Sprintf("credit:used.free_%s", userId)).Int64()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, err
		}

		total := helper.Select[int64](strings.HasSuffix(userId, "1"), 0, 0)
		if creditUsed > total {
			creditUsed = total
		}

		return &paymentpb.CreditState{
			UserId:    userId,
			Total:     total,
			Used:      creditUsed,
			PlanId:    "free",
			PaymentId: "free",
		}, nil

	} else {

		plan := configs.GetPlanById(ongoingPayment.PlanID)

		creditUsed, err := t.data.Repos.RedisClient.Get(ctx, fmt.Sprintf("credit:used.%d_%s", ongoingPayment.ID, userId)).Int64()
		if err != nil && !errors.Is(err, redis.Nil) {
			return nil, err
		}

		if creditUsed > plan.CreditPerMonth {
			creditUsed = plan.CreditPerMonth
		}

		remaining := plan.CreditPerMonth - creditUsed
		if remaining < 0 {
			remaining = 0
		}

		return &paymentpb.CreditState{
			UserId:    userId,
			Total:     plan.CreditPerMonth,
			Used:      creditUsed,
			Remaining: remaining,
			PlanId:    plan.Id,
			PaymentId: conv.String(ongoingPayment.ID),
			ExpireAt:  ongoingPayment.ExpireAt.Unix(),
		}, nil
	}

}

func (t *PaymentBiz) GetOngoingPayment(ctx context.Context, userId string) (*ent.Payment, error) {

	cachedPayment, exists := t.data.Repos.LocalCache.Get("ongoingPayment:" + userId)
	if exists {

		_, ok := cachedPayment.(string)

		if ok {
			return nil, nil
		}

		return cachedPayment.(*ent.Payment), nil
	}

	ongoingPayments, err := t.data.Repos.EntClient.Payment.Query().
		Where(payment.UserIDEQ(userId)).
		Where(payment.ExpireAtGT(time.Now())).
		Where(payment.StatusEQ(enums.PaymentStatus_Complete)).
		Order(ent.Desc(payment.FieldExpireAt)).
		All(ctx)

	if err != nil {
		return nil, err
	}

	if len(ongoingPayments) == 0 {
		t.data.Repos.LocalCache.Set("ongoingPayment:"+userId, "null", time.Minute*5)
		return nil, nil
	}

	t.data.Repos.LocalCache.Set("ongoingPayment:"+userId, ongoingPayments[0], time.Minute*5)

	return ongoingPayments[0], nil
}

func (t *PaymentBiz) CreateV2(ctx context.Context, userId string, params *paymentpb.CreateParams) (string, error) {

	user, err := t.data.GrpcClients.UserClient.GetUserById(ctx, &userpb.GetUserByIdParams{Id: userId})
	if err != nil {
		return "", err
	}

	client := t.data.StripeFactory.Get()
	// 查找price
	lookupKey := fmt.Sprintf("%s-%s", params.Level, params.Cycle)

	priceIter := client.V1Prices.List(ctx, &stripe.PriceListParams{
		LookupKeys: stripe.StringSlice([]string{lookupKey}),
		Expand:     stripe.StringSlice([]string{"data.product"}),
	})

	var price *stripe.Price
	for x := range priceIter {
		price = x
		break
	}

	if price == nil {
		return "", errors.BadRequest(fmt.Sprintf("no price found: by %s, %s", params.Level, params.Cycle), "")
	}

	plan := configs.GetPlanById(lookupKey)

	pay, err := t.data.Repos.EntClient.Payment.Create().
		SetStatus(enums.PaymentStatus_Created).
		SetPlatform("stripe").
		SetAmount(plan.CnyAmount).
		SetPlanID(lookupKey).
		SetUserID(user.Id).
		//OnConflictColumns(payment.FieldPlatform, payment.FieldSubID).
		//UpdateExpireAt().
		Save(ctx)
	if err != nil {
		return "", err
	}

	//subscriptionData := &stripe.CheckoutSessionCreateSubscriptionDataParams{}
	//// 试用
	//if strings.HasPrefix(params.Cycle, "tried") || strings.HasPrefix(params.Cycle, "trial") {
	//	days := strings.Split(params.Cycle, "_")[0][5:]
	//	subscriptionData.TrialSettings = &stripe.CheckoutSessionSubscriptionDataTrialSettingsParams{
	//		EndBehavior: &stripe.CheckoutSessionSubscriptionDataTrialSettingsEndBehaviorParams{
	//			MissingPaymentMethod: stripe.String("cancel"),
	//		},
	//	}
	//	subscriptionData.TrialPeriodDays = stripe.Int64(conv.Int64(days))
	//}
	//// 折扣
	//var discounts []*stripe.CheckoutSessionCreateDiscountParams
	//
	//if params.CouponID != "" {
	//	coupon, err := client.V1Coupons.Retrieve(ctx, params.CouponID, nil)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	if coupon != nil && coupon.Valid {
	//		discounts = append(discounts, &stripe.CheckoutSessionCreateDiscountParams{
	//			Coupon: stripe.String(params.CouponID),
	//		})
	//	}
	//}
	//
	////
	//if params.PromotionCode != "" {
	//	promIter := client.V1PromotionCodes.List(ctx, &stripe.PromotionCodeListParams{
	//		Code: stripe.String(params.PromotionCode),
	//	})
	//
	//	for x := range promIter {
	//		log.Debugw("promotionCode", x)
	//
	//		if x.Active {
	//			discounts = append(discounts, &stripe.CheckoutSessionCreateDiscountParams{
	//				Coupon: stripe.String(x.Coupon.ID),
	//			})
	//		}
	//	}
	//
	//}

	var customFields []*stripe.CheckoutSessionCreateCustomFieldParams
	if user.Phone == "" {
		customFields = append(customFields, &stripe.CheckoutSessionCreateCustomFieldParams{
			Key:      stripe.String("phoneNumber"), // 用于识别该字段的唯一键
			Type:     stripe.String("text"),        // 使用文本类型
			Optional: stripe.Bool(false),           // 设置为必填字段
			Label: &stripe.CheckoutSessionCreateCustomFieldLabelParams{
				Type:   stripe.String("custom"),
				Custom: stripe.String("手机号码"),
			},
			Text: &stripe.CheckoutSessionCreateCustomFieldTextParams{
				MinimumLength: stripe.Int64(11), // 中国手机号是11位
				MaximumLength: stripe.Int64(11), // 限制最大长度也是11位
				// 使用正则表达式验证手机号格式
				//Pattern:   stripe.String("^1[3-9]\\d{9}$"),
			},
		},
		)
	}

	var email *string
	if user.Email != "" {
		email = &user.Email
	}

	session, err := client.V1CheckoutSessions.Create(ctx, &stripe.CheckoutSessionCreateParams{
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{Price: &price.ID, Quantity: stripe.Int64(1)},
		},
		Mode: stripe.String(string(stripe.CheckoutSessionModePayment)),
		//Discounts:                discounts,
		AllowPromotionCodes:      stripe.Bool(true),
		SuccessURL:               &params.SuccessUrl,
		CancelURL:                &params.CancelUrl,
		ClientReferenceID:        stripe.String(conv.String(pay.ID)),
		CustomerEmail:            email,
		CustomerCreation:         stripe.String("if_required"),
		BillingAddressCollection: stripe.String("auto"),
		//Customer:          &user.Email,
		// 添加自定义表单
		CustomFields: customFields,
		//PaymentMethodCollection: stripe.String(string(stripe.CheckoutSessionPaymentMethodCollectionIfRequired)),
		//SubscriptionData: subscriptionData,
		PaymentMethodTypes: stripe.StringSlice([]string{
			string(stripe.PaymentMethodTypeWeChatPay),
			string(stripe.PaymentMethodTypeAlipay),
			string(stripe.PayoutDestinationTypeCard),
		}),
		PaymentMethodOptions: &stripe.CheckoutSessionCreatePaymentMethodOptionsParams{
			WeChatPay: &stripe.CheckoutSessionCreatePaymentMethodOptionsWeChatPayParams{
				Client: stripe.String("web"),
			},
		},
	})

	if err != nil {
		return "", err
	}

	log.Debugw("session ", conv.S2J(session))

	return session.URL, nil
}

// 辅助函数：获取或创建客户
func (t *PaymentBiz) getOrCreateCustomer(ctx context.Context, email string) (string, error) {
	// 首先查找现有客户
	client := t.data.StripeFactory.Get()

	params := &stripe.CustomerListParams{
		Email: stripe.String(email),
		ListParams: stripe.ListParams{
			Limit: stripe.Int64(1),
		},
	}

	iter := client.V1Customers.List(ctx, params)

	for x := range iter {
		return x.ID, nil
	}

	// 创建新客户
	cus, err := client.V1Customers.Create(ctx, &stripe.CustomerCreateParams{
		Params: stripe.Params{},
		Email:  stripe.String(email),
		//InvoiceSettings: &stripe.CustomerCreateInvoiceSettingsParams{
		//	DefaultPaymentMethod: stripe.String(string(stripe.PaymentMethodTypeAlipay)),
		//},
		Metadata: map[string]string{
			//"created_from": "mobile_app",
		},
		//PaymentMethod: stripe.String(string(stripe.PaymentMethodTypeAlipay)),
	})
	if err != nil {
		return "", err
	}

	return cus.ID, nil
}

func (t *PaymentBiz) CreateSubscription(ctx context.Context, userId string, params *paymentpb.CreateParams) (string, error) {

	client := t.data.StripeFactory.Get()

	// 查找price
	lookupKey := fmt.Sprintf("%s-%s", params.Level, params.Cycle)

	priceIter := client.V1Prices.List(ctx, &stripe.PriceListParams{
		LookupKeys: stripe.StringSlice([]string{lookupKey}),
		Expand:     stripe.StringSlice([]string{"data.product"}),
	})

	var price *stripe.Price
	for x := range priceIter {
		price = x
		break
	}

	if price == nil {
		return "", errors.BadRequest(fmt.Sprintf("no price found: by %s, %s", params.Level, params.Cycle), "")
	}

	customerId, err := t.getOrCreateCustomer(ctx, "tongxurt@gmail.com")
	if err != nil {
		return "", err
	}
	//
	//pmIter := client.V1PaymentMethods.List(ctx, &stripe.PaymentMethodListParams{
	//	Customer: stripe.String(customerId),
	//})
	//
	//for x := range pmIter {
	//	fmt.Println(x)
	//}

	//key, err := client.V1EphemeralKeys.Create(ctx, &stripe.EphemeralKeyCreateParams{
	//	Customer: stripe.String(customerId),
	//})
	//if err != nil {
	//	return "", err
	//}
	//
	//fmt.Println(key)
	//
	//create, err := client.V1SetupIntents.Create(ctx, &stripe.SetupIntentCreateParams{
	//	Params:       stripe.Params{},
	//	AttachToSelf: nil,
	//	AutomaticPaymentMethods: &stripe.SetupIntentCreateAutomaticPaymentMethodsParams{
	//		AllowRedirects: stripe.String("always"),
	//		Enabled:        stripe.Bool(true),
	//	},
	//	Confirm:                    nil,
	//	ConfirmationToken:          nil,
	//	Customer:                   stripe.String(customerId),
	//	Description:                nil,
	//	Expand:                     nil,
	//	FlowDirections:             nil,
	//	MandateData:                nil,
	//	Metadata:                   nil,
	//	OnBehalfOf:                 nil,
	//	PaymentMethod:              nil,
	//	PaymentMethodConfiguration: nil,
	//	PaymentMethodData:          nil,
	//	PaymentMethodOptions:       &stripe.SetupIntentCreatePaymentMethodOptionsParams{},
	//	//PaymentMethodTypes: []*string{
	//	//	stripe.String(stripe.PaymentMethodTypeAlipay),
	//	//	stripe.String(stripe.PaymentMethodTypeCard),
	//	//},
	//	ReturnURL:    nil,
	//	SingleUse:    nil,
	//	Usage:        stripe.String("off_session"),
	//	UseStripeSDK: stripe.Bool(true),
	//})
	//if err != nil {
	//	return "", err
	//}

	create, err := client.V1PaymentIntents.Create(ctx, &stripe.PaymentIntentCreateParams{
		Params:               stripe.Params{},
		Amount:               stripe.Int64(price.UnitAmount),
		ApplicationFeeAmount: nil,
		//AutomaticPaymentMethods: &stripe.PaymentIntentCreateAutomaticPaymentMethodsParams{
		//	AllowRedirects: stripe.String("always"),
		//	Enabled:        stripe.Bool(true),
		//},
		PaymentMethodTypes: []*string{
			stripe.String(stripe.PaymentMethodTypeAlipay),
			stripe.String(stripe.PaymentMethodTypeCard),
		},
		CaptureMethod: nil,
		//Confirm:               stripe.Bool(true),
		ConfirmationMethod:    nil,
		ConfirmationToken:     nil,
		Currency:              stripe.String(price.Currency),
		Customer:              stripe.String(customerId),
		Description:           nil,
		ErrorOnRequiresAction: nil,
		Expand:                nil,
		Mandate:               nil,
		MandateData:           nil,
		Metadata:              nil,
		//OffSession:            stripe.Bool(true),

		SetupFutureUsage: stripe.String("off_session"),
		UseStripeSDK:     nil,
	})
	if err != nil {
		return "", err
	}

	return create.ClientSecret, nil

	//create, err := client.V1Subscriptions.Create(ctx, &stripe.SubscriptionCreateParams{
	//	Params: stripe.Params{
	//		//Expand: stripe.StringSlice([]string{"latest_invoice.payment_intent"}),
	//	},
	//	AddInvoiceItems:             nil,
	//	ApplicationFeePercent:       nil,
	//	AutomaticTax:                nil,
	//	BackdateStartDate:           nil,
	//	BillingCycleAnchor:          nil,
	//	BillingCycleAnchorConfig:    nil,
	//	BillingCycleAnchorNow:       nil,
	//	BillingCycleAnchorUnchanged: nil,
	//	CancelAt:                    nil,
	//	CancelAtPeriodEnd:           nil,
	//	CollectionMethod:            nil,
	//	Currency:                    nil,
	//	Customer:                    stripe.String(customerId),
	//	//CustomerEmail:            email,
	//	DaysUntilDue: nil,
	//	//DefaultPaymentMethod: stripe.String(string(stripe.PaymentMethodTypeAlipay)),
	//	DefaultSource:   nil,
	//	DefaultTaxRates: nil,
	//	Description:     nil,
	//	Discounts:       nil,
	//	Expand:          nil,
	//	InvoiceSettings: nil,
	//	Items: []*stripe.SubscriptionCreateItemParams{
	//		{Price: &price.ID, Quantity: stripe.Int64(1)},
	//	},
	//
	//	Metadata:        nil,
	//	OffSession:      nil,
	//	OnBehalfOf:      nil,
	//	PaymentBehavior: stripe.String("default_incomplete"),
	//	PaymentSettings: &stripe.SubscriptionCreatePaymentSettingsParams{
	//		SaveDefaultPaymentMethod: stripe.String("on_subscription"),
	//	},
	//	PendingInvoiceItemInterval: nil,
	//	ProrationBehavior:          nil,
	//	TransferData:               nil,
	//	TrialEnd:                   nil,
	//	TrialEndNow:                nil,
	//	TrialFromPlan:              nil,
	//	TrialPeriodDays:            nil,
	//	TrialSettings:              nil,
	//})
	//if err != nil {
	//	return "", err
	//}
	//
	//return create.LatestInvoice.ConfirmationSecret.ClientSecret, nil
}

func (t *PaymentBiz) Create(ctx context.Context, userId string, params *paymentpb.CreateParams) (string, error) {

	user, err := t.data.GrpcClients.UserClient.GetUserById(ctx, &userpb.GetUserByIdParams{Id: userId})
	if err != nil {
		return "", err
	}

	client := t.data.StripeFactory.Get()
	// 查找price
	lookupKey := fmt.Sprintf("%s-%s", params.Level, params.Cycle)

	priceIter := client.V1Prices.List(ctx, &stripe.PriceListParams{
		LookupKeys: stripe.StringSlice([]string{lookupKey}),
		Expand:     stripe.StringSlice([]string{"data.product"}),
	})

	var price *stripe.Price
	for x := range priceIter {
		price = x
		break
	}

	if price == nil {
		return "", errors.BadRequest(fmt.Sprintf("no price found: by %s, %s", params.Level, params.Cycle), "")
	}

	subscriptionData := &stripe.CheckoutSessionCreateSubscriptionDataParams{}
	//// 试用
	//if strings.HasPrefix(params.Cycle, "tried") || strings.HasPrefix(params.Cycle, "trial") {
	//	days := strings.Split(params.Cycle, "_")[0][5:]
	//	subscriptionData.TrialSettings = &stripe.CheckoutSessionSubscriptionDataTrialSettingsParams{
	//		EndBehavior: &stripe.CheckoutSessionSubscriptionDataTrialSettingsEndBehaviorParams{
	//			MissingPaymentMethod: stripe.String("cancel"),
	//		},
	//	}
	//	subscriptionData.TrialPeriodDays = stripe.Int64(conv.Int64(days))
	//}
	// 折扣
	var discounts []*stripe.CheckoutSessionCreateDiscountParams

	if params.CouponID != "" {
		coupon, err := client.V1Coupons.Retrieve(ctx, params.CouponID, nil)
		if err != nil {
			return "", err
		}

		if coupon != nil && coupon.Valid {
			discounts = append(discounts, &stripe.CheckoutSessionCreateDiscountParams{
				Coupon: stripe.String(params.CouponID),
			})
		}
	}

	//
	if params.PromotionCode != "" {
		promIter := client.V1PromotionCodes.List(ctx, &stripe.PromotionCodeListParams{
			Code: stripe.String(params.PromotionCode),
		})

		for x := range promIter {
			log.Debugw("promotionCode", x)

			if x.Active {
				discounts = append(discounts, &stripe.CheckoutSessionCreateDiscountParams{
					Coupon: stripe.String(x.Coupon.ID),
				})
			}
		}

	}

	session, err := client.V1CheckoutSessions.Create(ctx, &stripe.CheckoutSessionCreateParams{
		LineItems: []*stripe.CheckoutSessionCreateLineItemParams{
			{Price: &price.ID, Quantity: stripe.Int64(1)},
		},
		Mode:              stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		Discounts:         discounts,
		SuccessURL:        &params.SuccessUrl,
		CancelURL:         &params.CancelUrl,
		ClientReferenceID: types.NewClientReference(userId, params.PromotionCode).ToString(),
		CustomerEmail:     &user.Email,
		//Customer:          &user.Email,
		// 添加自定义表单
		CustomFields: []*stripe.CheckoutSessionCreateCustomFieldParams{{}},
		//PaymentMethodCollection: stripe.String(string(stripe.CheckoutSessionPaymentMethodCollectionIfRequired)),
		SubscriptionData: subscriptionData,
	})

	if err != nil {
		return "", err
	}

	log.Debugw("session ", conv.S2J(session))

	return session.URL, nil
}

func (t *PaymentBiz) OnEvent(ctx context.Context, body []byte, signature, callback string) error {

	client := t.data.StripeFactory.GetByCallback(callback)

	callbackEvent, err := client.ConstructEvent(body, signature)
	if err != nil {
		return err
	}
	log.Debugw("Stripe OnEvent ", "", "callbackEvent", callbackEvent.Type, "body", string(body))

	switch callbackEvent.Type {
	case stripe.EventTypeInvoicePaid:
		invoiceID := conv.String(callbackEvent.Data.Object["id"])
		err = t.onPaid(ctx, client, invoiceID)
	case stripe.EventTypeCheckoutSessionCompleted:
		clientReferenceId := conv.String(callbackEvent.Data.Object["client_reference_id"])

		err = t.data.Repos.EntClient.Payment.Update().
			SetStatus(enums.PaymentStatus_Complete).
			SetExpireAt(time.Now().Add(24 * time.Hour * 30)).
			Where(payment.ID(conv.Int64(clientReferenceId))).
			//OnConflictColumns(payment.FieldPlatform, payment.FieldSubID).
			//UpdateExpireAt().
			Exec(ctx)
		if err != nil {
			return err
		}

		// todo
		t.data.Repos.LocalCache.Flush()

	default:
	}

	return err
}

func (t *PaymentBiz) onPaid(ctx context.Context, client *stripefactory.Client, invoiceID string) error {

	invoice, err := client.V1Invoices.Retrieve(ctx, invoiceID, nil)
	if err != nil {
		return err
	}

	subID := invoice.Parent.SubscriptionDetails.Subscription.ID
	sub, err := client.V1Subscriptions.Retrieve(ctx, subID, nil)
	if err != nil {
		return err
	}

	user, err := t.data.GrpcClients.UserClient.GetUserByEmail(ctx, &userpb.GetUserByEmailParams{
		Email: invoice.CustomerEmail,
	})
	if err != nil {
		return err
	}

	expireAt := sub.Items.Data[0].CurrentPeriodEnd
	planId := sub.Items.Data[0].Price.LookupKey
	amount := sub.Items.Data[0].Plan.Amount

	err = t.data.Repos.EntClient.Payment.Create().
		SetStatus(enums.PaymentStatus_Complete).
		SetSubID(subID).
		SetPlatform("stripe").
		SetAmount(float64(amount)).
		SetExpireAt(time.Unix(expireAt, 0)).
		SetPlanID(planId).
		SetUserID(user.Id).
		SetSessionID(invoice.ID).
		//OnConflictColumns(payment.FieldPlatform, payment.FieldSubID).
		//UpdateExpireAt().
		Exec(ctx)
	if err != nil {
		return err
	}

	t.data.Repos.LocalCache.Delete("ongoingPayment:" + user.Id)

	return nil
}

func (t *PaymentBiz) GetStripeBillUrl(ctx context.Context, version, email string) (string, error) {

	//client := t.data.StripeFactory.GetByEmail(email)
	//
	//sc := client.V1Customers.List(ctx, &stripe.CustomerListParams{
	//	Email: stripe.String(email),
	//})
	//
	//for _, x := range sc {
	//
	//}
	//
	//if iter.Err() != nil {
	//	return "", iter.Err()
	//}
	//
	//iter
	//
	//if len(iter.CustomerList().Data) == 0 {
	//	return client.BillLoginUrl(), nil
	//}
	//
	//customer := iter.CustomerList().Data[0]
	//
	//params := &stripe.BillingPortalSessionParams{
	//	//Configuration: stripe.String(conv.S2J(configuration)),
	//	Customer: stripe.String(customer.ID),
	//	//ReturnURL: stripe.String(t.data.Conf.Component.Stripe.BillLoginUrl),
	//}
	//
	//session, err := client.BillingPortalSessions.New(params)
	//if err != nil {
	//	return "", err
	//}
	//
	//return session.URL, nil

	return "", nil
}
