package stripe

import (
	stripeClient "github.com/stripe/stripe-go/v82"
)

type Config struct {
	Key          string
	Secret       string
	BillLoginUrl string
	Emails       []string
}

//func NewStripeClient(conf Config) *stripeClient.API {
//	return stripeClient.New(conf.Field, nil)
//}

type Client struct {
	*stripeClient.Client
}

func NewStripeClient(conf Config) *Client {
	return &Client{
		Client: stripeClient.NewClient(conf.Key, nil),
	}
}
