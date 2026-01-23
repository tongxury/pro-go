package tikhub

import (
	"context"
	"testing"
)

func TestClient_XhsSearchUsers(t *testing.T) {

	c := NewClient()
	//
	//c.xhsAppSearchUsers(context.Background(), "6769882234")
	c.XhsSearchUsers(context.Background(), "大润发")
}
