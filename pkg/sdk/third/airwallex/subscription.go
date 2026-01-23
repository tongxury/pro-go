package airwallex

import "context"

func (t *Client) CreateSubscription(ctx context.Context, params CreateSubscriptionParams) (*ProductList, error) {

	var r ProductList

	err := t.post(ctx, "/api/v1/subscriptions/create",
		map[string]interface{}{
			"billing_customer_id":       params.BillingCustomerId,
			"collection_method":         "AUTO_CHARGE",
			"linked_payment_account_id": t.config.AccountId,
			"payment_source_id":         "",
			"days_until_due":            0,
			"default_invoice_template": map[string]interface{}{
				"invoice_memo": "Thanks for your purchase.",
			},
			"items":      params.Items,
			"recurring":  params.Recurring,
			"request_id": params.RequestId,
		}, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type Item struct {
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	PriceId  string                 `json:"price_id,omitempty"`
	Quantity int                    `json:"quantity,omitempty"`
}

type Recurring struct {
	Period     int    `json:"period"`
	PeriodUnit string `json:"period_unit"`
}
type CreateSubscriptionParams struct {
	BillingCustomerId string                 `json:"billing_customer_id,omitempty"`
	Items             []Item                 `json:"items"`
	Metadata          map[string]interface{} `json:"metadata"`
	Recurring         Recurring              `json:"recurring"`
	RequestId         string                 `json:"request_id"`
}
