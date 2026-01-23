package gmgnai

import (
	"context"
	"fmt"
	"testing"
)

func TestGetLatestTrades(t *testing.T) {

	c := NewClient()

	trades, err := c.GetLatestTrades(context.Background(), GetLatestTradesParams{
		Chain: "sol",
		Token: "4jF8K9WZ4NVto9gtmCHaV8rQb2Z6UhJYQqSZNhyipump",
		Limit: 100,
	})
	if err != nil {
		return
	}

	fmt.Println(trades)
}
