package sdk

import (
	"context"
	"time"
)

type ICache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, dur time.Duration) error
}
