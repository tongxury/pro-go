package wxpay

import (
	"context"
	"testing"
)

func TestClient_ParseWebhookData(t *testing.T) {

	c := NewClient()

	c.ParseWebhookData(context.Background())
}
