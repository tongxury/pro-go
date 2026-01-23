package bitquery

import (
	"context"
	"testing"
)

func TestFindLatestTrades(t *testing.T) {

	c := New()

	_, _ = c.FindLatestTrades(context.Background(), FindLatestTradesParams{
		Token: "5SAbu2zKEuunG6QWYpqaMDnKG33yFWJe2V2MuyG1pump",
		Limit: 30,
	})

}
