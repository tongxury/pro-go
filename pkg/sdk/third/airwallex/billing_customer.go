package airwallex

import (
	"context"
	"fmt"
	"github.com/google/uuid"
)

func (t *Client) GetBillingCustomer(ctx context.Context, customerId string) (*BillingCustomer, error) {

	var r BillingCustomer
	err := t.get(ctx, fmt.Sprintf("/api/v1/billing_customers/%s", customerId), nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (t *Client) CreateBillingCustomer(ctx context.Context, params CreateBillingCustomerParams) (*BillingCustomer, error) {

	var r BillingCustomer

	err := t.post(ctx, "/api/v1/billing_customers/create", map[string]interface{}{
		"default_legal_entity_id": "le_qv8IH-4cPsG8XCh095pAnQ",
		"request_id":              uuid.NewString(),
		"email":                   params.Email,
		"name":                    params.Name,
		"nickname":                params.Nickname,
		"phone_number":            params.PhoneNumber,
		"metadata":                params.Metadata,
	}, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type CreateBillingCustomerParams struct {
	//DefaultLegalEntityId string
	Email       string                 `json:"email"`
	Name        string                 `json:"name"`
	Nickname    string                 `json:"nickname"`
	PhoneNumber string                 `json:"phone_number"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type BillingCustomer struct {
	Address                 Address                `json:"address"`
	DefaultBillingCurrency  string                 `json:"default_billing_currency"`
	DefaultLegalEntityId    string                 `json:"default_legal_entity_id"`
	Description             string                 `json:"description"`
	Email                   string                 `json:"email"`
	Metadata                map[string]interface{} `json:"metadata"`
	Name                    string                 `json:"name"`
	Nickname                string                 `json:"nickname"`
	PhoneNumber             string                 `json:"phone_number"`
	RequestId               string                 `json:"request_id"`
	TaxIdentificationNumber string                 `json:"tax_identification_number"`
	Type                    string                 `json:"type"`
}
