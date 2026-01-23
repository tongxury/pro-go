package helper

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRetry(t *testing.T) {

	fn1 := func(ctx context.Context) error {
		return errors.New("error")
	}

	ctx := context.Background()
	times, err := Retry(ctx, fn1, 10)
	assert.Equal(t, 10, times)
	assert.Equal(t, err.Error(), "error")

	fn2 := func(ctx context.Context) error {
		return nil
	}

	times2, err2 := Retry(ctx, fn2, 10)
	assert.Equal(t, 0, times2)
	assert.Nil(t, err2)
}
