package volcengine

import (
	"context"
	"testing"
)

func TestClient_GetFlowerList(t *testing.T) {

	c := NewClient()

	c.GetFlowerList(
		context.Background(),
		GetFlowerListParams{},
	)
}
