package volcengine

import (
	"context"
	"testing"
)

func TestClient_ListUsers(t *testing.T) {

	c := NewClient()

	c.ListUsers(
		context.Background(),
	)
}
