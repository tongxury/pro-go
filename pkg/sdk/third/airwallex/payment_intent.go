package airwallex

import (
	"context"
	"github.com/google/uuid"
)

func (t *Client) CreatePaymentIntent(ctx context.Context, params CreatePaymentIntentParams) (*PaymentIntent, error) {

	var r PaymentIntent

	err := t.post(ctx, "/api/v1/pa/payment_intents/create", map[string]interface{}{
		"request_id":        uuid.NewString(),
		"merchant_order_id": params.MerchantOrderId,
		"amount":            params.Product.UnitPrice * float64(params.Product.Quantity),
		"currency":          params.Currency,
		"customer_id":       params.CustomerId,
		"order": map[string]interface{}{
			"products": []map[string]interface{}{
				{
					"name":         params.Product.Name,
					"description":  params.Product.Description,
					"unit_price: ": params.Product.UnitPrice,
					"quantity":     params.Product.Quantity,
				},
			},
		},
		"return_url": params.ReturnUrl,
		"metadata":   params.Metadata,
	}, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type Order struct {
	Products []Product `json:"products"`
}
type PaymentIntent struct {
	Id                          string   `json:"id"`
	RequestId                   string   `json:"request_id"`
	Amount                      int      `json:"amount"`
	Currency                    string   `json:"currency"`
	MerchantOrderId             string   `json:"merchant_order_id"`
	Order                       Order    `json:"order"`
	Descriptor                  string   `json:"descriptor"`
	Status                      string   `json:"status"`
	CapturedAmount              int      `json:"captured_amount"`
	CreatedAt                   string   `json:"created_at"`
	UpdatedAt                   string   `json:"updated_at"`
	AvailablePaymentMethodTypes []string `json:"available_payment_method_types"`
	ClientSecret                string   `json:"client_secret"`
	BaseAmount                  int      `json:"base_amount"`
	BaseCurrency                string   `json:"base_currency"`
}

type CreatePaymentIntentParams struct {
	Amount          int64
	MerchantOrderId string
	Product         CreatePaymentIntentParams_Product
	CustomerId      string
	Currency        string
	ReturnUrl       string
	Metadata        map[string]any
}

type CreatePaymentIntentParams_Product struct {
	//Code           string  `json:"code"`
	//CommodityCode  string  `json:"commodity_code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	//DiscountAmount float64 `json:"discount_amount"`
	Quantity int `json:"quantity"`
	//TaxPercent     float64 `json:"tax_percent"`
	//TotalAmount    float64 `json:"total_amount"`
	//TotalTaxAmount float64 `json:"total_tax_amount"`
	Unit      string  `json:"unit"`
	UnitPrice float64 `json:"unit_price"`
}

type T struct {
	RequestId       string `json:"request_id"`
	Amount          int    `json:"amount"`
	Currency        string `json:"currency"`
	MerchantOrderId string `json:"merchant_order_id"`
	Metadata        struct {
		MyTestMetadataId string `json:"my_test_metadata_id"`
	} `json:"metadata"`
	Order struct {
		Products []struct {
			Code      string  `json:"code"`
			Desc      string  `json:"desc"`
			Name      string  `json:"name"`
			Quantity  int     `json:"quantity"`
			Sku       string  `json:"sku"`
			Type      string  `json:"type"`
			UnitPrice float64 `json:"unit_price"`
			Url       string  `json:"url"`
		} `json:"products"`
		Shipping struct {
			Address struct {
				City        string `json:"city"`
				CountryCode string `json:"country_code"`
				Postcode    string `json:"postcode"`
				State       string `json:"state"`
				Street      string `json:"street"`
			} `json:"address"`
			FirstName      string `json:"first_name"`
			LastName       string `json:"last_name"`
			PhoneNumber    string `json:"phone_number"`
			ShippingMethod string `json:"shipping_method"`
		} `json:"shipping"`
		Type string `json:"type"`
	} `json:"order"`
	PaymentMethodOptions struct {
		Card struct {
			RiskControl struct {
				SkipRiskProcessing      bool   `json:"skip_risk_processing"`
				ThreeDomainSecureAction string `json:"three_domain_secure_action"`
				ThreeDsAction           string `json:"three_ds_action"`
			} `json:"risk_control"`
		} `json:"card"`
	} `json:"payment_method_options"`
	ReturnUrl string `json:"return_url"`
}
