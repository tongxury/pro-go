package tikhub

import (
	"context"
	"fmt"
	"testing"
)

func TestClient_XhsGetUserById(t *testing.T) {

	c := NewClient()
	id, err := c.XhsGetUserById(context.Background(), "5f0d02c5000000000101e19e")
	if err != nil {
		return
	}

	fmt.Println(id)
}
