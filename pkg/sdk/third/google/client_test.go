package google

import (
	"context"
	"testing"
)

func TestName(t *testing.T) {
	c := NewClient()

	c.SendEmail(context.Background())
}
