package stripe

import (
	"context"
	"store/confs"
	"testing"

	"github.com/stripe/stripe-go/v82"
)

func TestNewStripeClient(t *testing.T) {

	c := NewStripeClient(Config{
		Key:    confs.StripeKeyTest,
		Secret: confs.StripeSecretTest,

		BillLoginUrl: "",
	})

	ctx := context.Background()

	c.V1CheckoutSessions.Create(ctx, &stripe.CheckoutSessionCreateParams{
		Params:                     stripe.Params{},
		AdaptivePricing:            nil,
		AfterExpiration:            nil,
		AllowPromotionCodes:        nil,
		AutomaticTax:               nil,
		BillingAddressCollection:   nil,
		CancelURL:                  nil,
		ClientReferenceID:          nil,
		ConsentCollection:          nil,
		Currency:                   nil,
		Customer:                   nil,
		CustomerCreation:           nil,
		CustomerEmail:              nil,
		CustomerUpdate:             nil,
		CustomFields:               nil,
		CustomText:                 nil,
		Discounts:                  nil,
		Expand:                     nil,
		ExpiresAt:                  nil,
		InvoiceCreation:            nil,
		LineItems:                  nil,
		Locale:                     nil,
		Metadata:                   nil,
		Mode:                       nil,
		OptionalItems:              nil,
		PaymentIntentData:          nil,
		PaymentMethodCollection:    nil,
		PaymentMethodConfiguration: nil,
		PaymentMethodData:          nil,
		PaymentMethodOptions:       nil,
		PaymentMethodTypes:         nil,
		Permissions:                nil,
		PhoneNumberCollection:      nil,
		RedirectOnCompletion:       nil,
		ReturnURL:                  nil,
		SavedPaymentMethodOptions:  nil,
		SetupIntentData:            nil,
		ShippingAddressCollection:  nil,
		ShippingOptions:            nil,
		SubmitType:                 nil,
		SubscriptionData:           nil,
		SuccessURL:                 nil,
		TaxIDCollection:            nil,
		UIMode:                     nil,
		WalletOptions:              nil,
	})
	//invoice, err := c.V1Invoices.Retrieve(ctx, "in_1RMmAAJ0m2xuUKf13gFb1JfV", nil)
	//if err != nil {
	//	t.Fatal(err)
	//}

	//// subscription object
	//sub, err := c.V1Subscriptions.Retrieve(ctx, invoice.Parent.SubscriptionDetails.Subscription.ID, nil)
	//if err != nil {
	//	t.Fatal(err)
	//}

}
