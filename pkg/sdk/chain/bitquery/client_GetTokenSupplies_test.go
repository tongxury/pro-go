package bitquery

import (
	"context"
	"testing"
)

func TestClientGetTokenSupplies(t *testing.T) {

	c := New()

	ctx := context.Background()

	_, _ = c.GetTokenSupplies(ctx, []string{"G9t6peRYVydgbjV5o9jvPCowcSrYhqABmnphZsQcpump"})

}
