package catalog

import (
	"fmt"
	"time"

	"context"

	"github.com/go-kit/kit/metrics"
)

type instrmw struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func InstrumentingMiddleware(counter metrics.Counter, latency metrics.Histogram) Middleware {
	return func(next Service) Service {
		return instrmw{
			requestCount:   counter,
			requestLatency: latency,
			next:           next,
		}
	}
}

func (mw instrmw) Search(ctx context.Context, query string) (books []Book, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "search", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	books, err = mw.next.Search(ctx, query)
	return
}

func (mw instrmw) Get(ctx context.Context, ID string) (book Book, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "login", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	book, err = mw.next.Get(ctx, ID)
	return
}
