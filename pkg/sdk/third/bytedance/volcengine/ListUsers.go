package volcengine

import (
	"context"
	"store/pkg/sdk/conv"
)

func (t *Client) ListUsers(ctx context.Context) {

	_, err := t.doRequest(ctx, Req{
		Version: "2022-02-01",
		Action:  "ListUsers",
		Service: "ic_iam",
		Method:  "POST",
		Body: conv.M2B(map[string]any{
			"UserType": "All",
		}),
	})
	if err != nil {
		return
	}
}
