package order

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

func (mw instrmw) PlaceOrder(ctx context.Context, bookID string) (order Order, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "place_order", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	order, err = mw.next.PlaceOrder(ctx, bookID)
	return
}

func (mw instrmw) GetUserOrders(ctx context.Context, userID string) (orders []Order, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "get_user_orders", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	orders, err = mw.next.GetUserOrders(ctx, userID)
	return
}

func (mw instrmw) CancelOrder(ctx context.Context, userID string, orderID string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "cancel_order", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.CancelOrder(ctx, userID, orderID)
	return
}
