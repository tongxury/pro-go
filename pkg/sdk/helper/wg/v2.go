package wg

import (
	"context"
	"store/pkg/sdk/helper"
	"sync"
)

func WaitGroupArgs[P any, T any](ctx context.Context, params []P, f func(ctx context.Context, x P) ([]T, error)) ([]T, error) {

	var errs []error
	var results []T

	wg := sync.WaitGroup{}
	l := sync.Mutex{}

	for i := range params {

		p := params[i]

		wg.Add(1)

		go func(ctx context.Context, param P) {
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

			if len(v) > 0 {
				l.Lock()
				results = append(results, v...)
				l.Unlock()
			}
		}(ctx, p)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return results, nil
}

func WaitGroupV2[T any](ctx context.Context, fs ...func(ctx context.Context) ([]T, error)) ([]T, error) {

	var errs []error
	var results []T

	wg := sync.WaitGroup{}
	l := sync.Mutex{}

	for i := range fs {

		f := fs[i]

		wg.Add(1)

		go func(ctx context.Context) {
			defer helper.DeferFunc(func() {
				wg.Done()
			})

			v, err := f(ctx)
			if err != nil {
				l.Lock()
				errs = append(errs, err)
				l.Unlock()
				return
			}

			if len(v) > 0 {
				l.Lock()
				results = append(results, v...)
				l.Unlock()
			}
		}(ctx)
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, errs[0]
	}

	return results, nil
}
