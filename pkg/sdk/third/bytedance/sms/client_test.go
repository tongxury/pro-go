package sms

import (
	"context"
	"testing"
)

func TestName(t *testing.T) {

	c := NewClient()

	c.SendSmsCode(context.Background(), "15701348086", "234124")
}
