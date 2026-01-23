package airwallex

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

func (t *Client) GetOrCreateCustomer(ctx context.Context, params CreateCustomerParams) (*Customer, error) {

	customers, err := t.ListCustomers(ctx, ListCustomersParams{
		MerchantCustomerId: params.MerchantCustomerId,
	})
	if err != nil {
		return nil, err
	}

	if customers != nil && len(customers.Items) > 0 {
		return &customers.Items[0], nil
	}

	customer, err := t.CreateCustomer(ctx, params)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (t *Client) GetCustomer(ctx context.Context, customerId string) (*Customer, error) {

	var r Customer
	err := t.get(ctx, fmt.Sprintf("/api/v1/pa/customers/%s", customerId), nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

func (t *Client) UpdateCustomer(ctx context.Context, params UpdateCustomerParams) (*Customer, error) {

	var r Customer

	err := t.post(ctx,
		fmt.Sprintf("/api/v1/pa/customers/%s/update", params.CustomerId),
		map[string]interface{}{
			"request_id":           uuid.NewString(),
			"merchant_customer_id": params.MerchantCustomerId,
			//"business_name": "",
			"first_name":   params.FirstName,
			"last_name":    params.LastName,
			"email":        params.Email,
			"phone_number": params.PhoneNumber,
			"metadata":     params.Metadata,
		}, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type ListCustomersParams struct {
	FromCreatedAt      *time.Time // ISO8601 format
	ToCreatedAt        *time.Time
	MerchantCustomerId string
	Page               int
	Size               int
}

func (t *Client) ListCustomers(ctx context.Context, params ListCustomersParams) (*List[Customer], error) {

	query := map[string]string{
		//"merchant_customer_id": params.MerchantCustomerId,
		//"page_num":             fmt.Sprintf("%d", params.Page),
		//"page_size":            fmt.Sprintf("%d", params.Size),
	}

	if params.MerchantCustomerId != "" {
		query["merchant_customer_id"] = params.MerchantCustomerId
	}

	if params.FromCreatedAt != nil {
		query["from_created_at"] = params.FromCreatedAt.Format("2006-01-02T15:04:05+0700")
	}

	if params.ToCreatedAt != nil {
		query["to_created_at"] = params.ToCreatedAt.Format("2006-01-02T15:04:05+0700")
	}

	if params.Page > 0 {
		query["page_num"] = fmt.Sprintf("%d", params.Page)
	}
	if params.Size > 0 {
		query["page_size"] = fmt.Sprintf("%d", params.Size)
	}

	var r List[Customer]
	err := t.get(ctx, "/api/v1/pa/customers", query, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil

}

func (t *Client) CreateCustomer(ctx context.Context, params CreateCustomerParams) (*Customer, error) {

	var r Customer

	err := t.post(ctx, "/api/v1/pa/customers/create", map[string]interface{}{
		"request_id":           uuid.NewString(),
		"merchant_customer_id": params.MerchantCustomerId,
		//"business_name": "",
		"first_name":   params.FirstName,
		"last_name":    params.LastName,
		"email":        params.Email,
		"phone_number": params.PhoneNumber,
		"metadata":     params.Metadata,
	}, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type UpdateCustomerParams struct {
	CustomerId string `json:"customer_id"`
	// 唯一
	MerchantCustomerId string                 `json:"merchant_customer_id"`
	FirstName          string                 `json:"first_name"`
	LastName           string                 `json:"last_name"`
	Email              string                 `json:"email"`
	PhoneNumber        string                 `json:"phone_number"`
	Metadata           map[string]interface{} `json:"metadata"`
}

type CreateCustomerParams struct {
	//RequestId          string `json:"request_id"`
	// 唯一
	MerchantCustomerId string                 `json:"merchant_customer_id"`
	FirstName          string                 `json:"first_name"`
	LastName           string                 `json:"last_name"`
	Email              string                 `json:"email"`
	PhoneNumber        string                 `json:"phone_number"`
	Metadata           map[string]interface{} `json:"metadata"`
}

type Address struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	Street      string `json:"street"`
	Postcode    string `json:"postcode"`
}
type Customer struct {
	Id                 string                 `json:"id"`
	RequestId          string                 `json:"request_id"`
	MerchantCustomerId string                 `json:"merchant_customer_id"`
	BusinessName       string                 `json:"business_name"`
	Email              string                 `json:"email"`
	PhoneNumber        string                 `json:"phone_number"`
	Address            Address                `json:"address"`
	Metadata           map[string]interface{} `json:"metadata"`
	ClientSecret       string                 `json:"client_secret"`
	CreatedAt          string                 `json:"created_at"`
	UpdatedAt          string                 `json:"updated_at"`
}
