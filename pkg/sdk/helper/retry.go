package helper

import (
	"context"
)

func Retry(ctx context.Context, fn func(ctx context.Context) error, maxRetries int) (int, error) {

	var times int
	for {

		err := fn(ctx)
		if err != nil {
			if times >= maxRetries {
				return times, err
			}
			times += 1

			continue
		}

		return times, nil
	}
}
