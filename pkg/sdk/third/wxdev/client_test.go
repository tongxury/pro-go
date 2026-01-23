package wxdev

import (
	"context"
	"testing"
)

func TestName(t *testing.T) {

	c := NewClient()

	c.GetUserPhoneNumber(context.Background(), "05a282ef27b0ad7ce70e14d10112398f33db4c181c1fe59f23353ce3b9c62edd")

}
