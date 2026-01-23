package service

import (
	"context"
	"fmt"
	"github.com/awa/go-iap/appstore"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"github.com/golang/protobuf/ptypes/empty"
	paymentpb "store/api/payment"
	"store/app/payment/configs"
)

func (t PaymentService) CallbackByAppleV2(ctx context.Context, params *paymentpb.CallbackByAppleV2Params) (*empty.Empty, error) {

	client := appstore.New()

	var tokenClaims jwt.MapClaims

	err := client.ParseNotificationV2WithClaim(params.SignedPayload, &tokenClaims)
	if err != nil {
		return nil, err
	}

	data, ok := tokenClaims["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid data")
	}

	var txClaims jwt.MapClaims

	err = client.ParseNotificationV2WithClaim(data["signedTransactionInfo"].(string), &txClaims)
	if err != nil {
		return nil, err
	}

	fmt.Println("appTransactionId", txClaims)
	//appTransactionId map[appTransactionId:704676719744237287 bundleId:com.tuturduck.veogoapp currency:CNY environment:Sandbox expiresDate:1.756321683e+12 inAppOwnershipType:PURCHASED originalPurchaseDate:1.754483949e+12 originalTransactionId:2000000977937351 price:38000 productId:veogo_l1_monthly purchaseDate:1.756321503e+12 quantity:1 signedDate:1.756321448792e+12 storefront:CHN storefrontId:143465 subscriptionGroupIdentifier:21718215 transactionId:2000000993966885 transactionReason:RENEWAL type:Auto-Renewable Subscription webOrderLineItemId:2000000109943837]
	//DEBUG CallbackByAppleV2= txClaims=map[appTransactionId:704676719744237287 bundleId:com.tuturduck.veogoapp currency:CNY environment:Sandbox expiresDate:1.756321683e+12 inAppOwnershipType:PURCHASED originalPurchaseDate:1.754483949e+12 originalTransactionId:2000000977937351 price:38000 productId:veogo_l1_monthly purchaseDate:1.756321503e+12 quantity:1 signedDate:1.756321448792e+12 storefront:CHN storefrontId:143465 subscriptionGroupIdentifier:21718215 transactionId:2000000993966885 transactionReason:RENEWAL type:Auto-Renewable Subscription webOrderLineItemId:2000000109943837]

	productId := txClaims["productId"].(string)
	//transactionId := txClaims["appTransactionId"].(string)

	//=map[appTransactionId:704800123819740572 bundleId:com.tuturduck.veogoapp currency:CNY environment:Production inAppOwnershipType:PURCHASED originalPurchaseDate:1.756311741e+12 originalTransactionId:460002625040184 price:38000 productId:com_veogo_l1_month purchaseDate:1.756311741e+12 quantity:1 signedDate:1.756311748329e+12 storefront:CHN storefrontId:143465 transactionId:460002625040184 transactionReason:PURCHASE type:Non-Consumable]

	//DEBUG CallbackByAppleV2= txClaims=map[appTransactionId:704676719744237287 bundleId:com.tuturduck.veogoapp currency:CNY environment:Sandbox expiresDate:1.756321683e+12 inAppOwnershipType:PURCHASED originalPurchaseDate:1.754483949e+12 originalTransactionId:2000000977937351 price:38000 productId:veogo_l1_monthly purchaseDate:1.756321503e+12 quantity:1 signedDate:1.756321448792e+12 storefront:CHN storefrontId:143465 subscriptionGroupIdentifier:21718215 transactionId:2000000993966885 transactionReason:RENEWAL type:Auto-Renewable Subscription webOrderLineItemId:2000000109943837]
	log.Debugw("CallbackByAppleV2", "", "txClaims", txClaims)

	//client.GetTransactionInfo()

	plan := configs.GetPlanById(productId)

	if plan == nil {
		return nil, errors.BadRequest("invalidPlan", "")
	}

	//err = t.data.Repos.EntClient.Payment.Create().
	//	SetStatus(enums.PaymentStatus_Complete).
	//	SetPlatform("apple").
	//	SetAmount(plan.Amount).
	//	SetExpireAt(time.Now().Add(time.Hour * 24 * 30)).
	//	SetPlanID(productId).
	//	SetSessionID(transactionId).
	//	SetUserID(userId).
	//	OnConflictColumns(payment.FieldSessionID).
	//	DoNothing().
	//	//UpdateExpireAt().
	//	Exec(ctx)
	//
	//if err != nil {
	//	log.Error("Payment Create Error", err, "userId", userId, "planId", plan.Id)
	//	return nil, err
	//}
	//
	//t.data.Repos.LocalCache.Delete("ongoingPayment:" + userId)

	return &empty.Empty{}, nil
}
