package stripefactory

import (
	stripego "github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	bizconf "store/app/payment/internal/conf"
	"store/pkg/confcenter"
	"store/pkg/sdk/third/stripe"
)

type Client struct {
	//version string
	conf   stripe.Config
	secret string
	*stripe.Client
}

type StripeFactory struct {
	prod *Client
	test *Client
	//prodV1 *Client
}

func NewStripeFactory(conf confcenter.StripePay) *StripeFactory {
	return &StripeFactory{
		prod: &Client{
			//version:      "",
			conf:   conf.Prod,
			secret: conf.Prod.Secret,
			Client: stripe.NewStripeClient(conf.Prod),
		},
		test: &Client{
			//version:      "",
			conf:   conf.Test,
			secret: conf.Test.Secret,
			Client: stripe.NewStripeClient(conf.Test),
		},
		//prodV1: &Client{
		//	conf:         conf.ProdV1,
		//	secret:       conf.ProdV1.Secret,
		//	Client: stripe.NewStripeClient(conf.ProdV1),
		//},
	}
}

func (t *Client) ConstructEvent(body []byte, signature string) (stripego.Event, error) {
	return webhook.ConstructEvent(body, signature, t.secret)
}

//func (t *Client) Version() string {
//	return t.version
//}

func (t *Client) BillLoginUrl() string {
	return t.conf.BillLoginUrl
}

//
//func (t *StripeFactory) Default() *Client {
//	return t.prod
//}
//func (t *StripeFactory) Test() *Client {
//	return t.test
//}

//func (t *StripeFactory) GetByEmail(email string) *Client {
//	if helper.InSlice(email, t.test.conf.Emails) {
//		return t.test
//	}
//
//	return t.prod
//}

func (t *StripeFactory) Get() *Client {
	return t.prod
}

//func (t *StripeFactory) GetByVersion(version string) *Client {
//	//if version == "v1" {
//	//	return t.prodV1
//	//}
//	return t.prod
//}

func (t *StripeFactory) GetByCallback(url string) *Client {
	switch url {
	case bizconf.StripeCallbackURL_Prod:
		return t.prod
	case bizconf.StripeCallbackURL_Test:
		return t.test
	//case bizconf.StripeCallbackURL_ProdV1:
	//	return t.prodV1
	default:
		return nil
	}
}
