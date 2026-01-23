package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/stripe/stripe-go/v82"
	paymentpb "store/api/payment"
	userpb "store/api/user"
	"store/app/payment/configs"
	"store/pkg/enums"
	"store/pkg/sdk/conv"
)

func (t *PaymentBiz) CreateV3(ctx context.Context, userId string, params *paymentpb.CreateParams) (string, error) {

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
		UIMode: stripe.String(string(stripe.CheckoutSessionUIModeCustom)),
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
