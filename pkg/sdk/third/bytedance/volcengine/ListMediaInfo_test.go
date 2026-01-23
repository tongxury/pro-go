package volcengine

import (
	"context"
	"testing"
)

func TestClient_ListMediaInfo(t *testing.T) {

	c := NewClient()

	c.ListMediaInfo(
		context.Background(),
		ListMediaInfoParams{},
	)
}
