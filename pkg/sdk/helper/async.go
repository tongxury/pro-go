package helper

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

func Go(ctx context.Context, f func(ctx context.Context)) {

	newCtx := context.Background()
	newCtx = opentracing.ContextWithSpan(newCtx, opentracing.SpanFromContext(ctx))

	go func(ctx2 context.Context) {

		defer DeferFunc()

		f(ctx2)

	}(newCtx)
}

func GoLoop(f func()) {
	go func() {
		defer DeferFunc()
		for {
			f()
		}
	}()
}

func GoArgs[T any](ctx context.Context, params T, f func(ctx context.Context, params T)) {

	newCtx := context.Background()
	newCtx = opentracing.ContextWithSpan(newCtx, opentracing.SpanFromContext(ctx))

	go func(ctx context.Context, params T) {

		defer DeferFunc()

		f(ctx, params)

	}(newCtx, params)

}
