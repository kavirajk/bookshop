package user

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

func (mw instrmw) Register(ctx context.Context, new NewUser) (user User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "register", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.next.Register(ctx, new)
	return
}

func (mw instrmw) Login(ctx context.Context, email, password string) (user User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "login", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.next.Login(ctx, email, password)
	return
}

func (mw instrmw) AuthToken(ctx context.Context, token string) (user User, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "auth_token", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	user, err = mw.next.AuthToken(ctx, token)
	return
}

func (mw instrmw) ResetPassword(ctx context.Context, key, newpass string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "reset_password", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.ResetPassword(ctx, key, newpass)
	return
}

func (mw instrmw) ChangePassword(ctx context.Context, userID string, oldpass, newpass string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "change_password", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.next.ChangePassword(ctx, userID, oldpass, newpass)
	return
}

func (mw instrmw) List(ctx context.Context, order string, limit, offset int) (users []User, total int, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "list", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	users, total, err = mw.next.List(ctx, order, limit, offset)
	return
}
