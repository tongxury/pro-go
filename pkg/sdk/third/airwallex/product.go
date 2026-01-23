package airwallex

import "context"

func (t *Client) ListProducts(ctx context.Context) (*List[Product], error) {

	var r List[Product]

	err := t.get(ctx, "/api/v1/products", nil, &r)
	if err != nil {
		return nil, err
	}

	return &r, nil
}

type Product struct {
	Active      bool                   `json:"active"`
	Description string                 `json:"description"`
	Id          string                 `json:"id"`
	Metadata    map[string]interface{} `json:"metadata"`
	Name        string                 `json:"name"`
	RequestId   string                 `json:"request_id"`
	Unit        string                 `json:"unit"`
}
type ProductList struct {
	HasMore bool      `json:"has_more"`
	Items   []Product `json:"items"`
}
