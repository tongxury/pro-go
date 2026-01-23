package bitquery

import (
	"context"
	"testing"
)

func TestClientGetChanges(t *testing.T) {

	c := New()

	_, _ = c.GetChanges(context.Background(), &GetChangesParams{
		Token:           "RBg96u3Z6GEfoigBLoN7pVBSKXgkEtPRyM8TW1Rpump",
		IntervalMinutes: 60,
		Limit:           30,
	})

}
