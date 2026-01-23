package bitquery

import (
	"context"
	"testing"
)

func TestGetTokenMetadataById(t *testing.T) {

	c := New()

	ctx := context.Background()

	_, _ = c.FindTokens(ctx, FindTokensParams{
		//Ids: []string{"G9t6peRYVydgbjV5o9jvPCowcSrYhqABmnphZsQcpump"},
		NameLike: "test",
	})

}
