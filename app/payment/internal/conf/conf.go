package bizconf

type BizConfig struct {
	//BaseUrl    string
	//SuccessUrl string
	//CancelUrl  string
	//Stripe clients.StripeConfig
}

const (
	StripeCallbackURL_Prod = "/api/pa/v1/pay/stripe-callback"
	StripeCallbackURL_Test = "/api/pa/v1/pay/stripe-callback-test"
	//StripeCallbackURL_ProdV1 = "/pay-api/stripe-callback"
)
