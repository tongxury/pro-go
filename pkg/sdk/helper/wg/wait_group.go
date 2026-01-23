package wg

import (
	"context"
	"store/pkg/sdk/helper"
	"sync"
)

func WaitGroupResults[T any, V any](ctx context.Context, params []T, f func(ctx context.Context, param T) (V, error)) ([]V, []error) {

	var errs []error
	var results []V

	wg := sync.WaitGroup{}
	l := sync.Mutex{}

	for i := range params {

		x := params[i]

		wg.Add(1)

		go func(ctx context.Context, param T) {
			defer helper.DeferFunc(func() {
				wg.Done()
			})

			v, err := f(ctx, param)
			if err != nil {
				l.Lock()
				errs = append(errs, err)
				l.Unlock()
				return
			}

			l.Lock()
			results = append(results, v)
			l.Unlock()

		}(ctx, x)
	}

	wg.Wait()

	return results, errs
}

func WaitGroupFunctions(ctx context.Context, fs ...func(ctx context.Context) error) []error {
	var errs []error

	wg := sync.WaitGroup{}
	l := sync.Mutex{}

	for i := range fs {

		x := fs[i]

		wg.Add(1)

		go func(ctx context.Context) {
			defer helper.DeferFunc(func() {
				wg.Done()
			})

			if err := x(ctx); err != nil {
				l.Lock()
				errs = append(errs, err)
				l.Unlock()
			}
		}(ctx)
	}

	wg.Wait()

	return errs
}

func WaitGroup[T any](ctx context.Context, params []T, f func(ctx context.Context, x T) error) []error {

	var errs []error

	wg := sync.WaitGroup{}
	l := sync.Mutex{}

	for i := range params {

		x := params[i]

		wg.Add(1)

		go func(ctx context.Context, param T) {
			defer helper.DeferFunc(func() {
				wg.Done()
			})

			if err := f(ctx, param); err != nil {
				l.Lock()
				errs = append(errs, err)
				l.Unlock()
			}
		}(ctx, x)
	}

	wg.Wait()

	return errs
}
func WaitGroupIndexed[T any](ctx context.Context, params []T, f func(ctx context.Context, x T, index int) error) []error {

	var errs []error

	wg := sync.WaitGroup{}
	l := sync.Mutex{}

	for i := range params {

		x := params[i]

		wg.Add(1)

		go func(ctx context.Context, param T) {
			defer helper.DeferFunc(func() {
				wg.Done()
			})

			if err := f(ctx, param, i); err != nil {
				l.Lock()
				errs = append(errs, err)
				l.Unlock()
			}
		}(ctx, x)
	}

	wg.Wait()

	return errs
}
