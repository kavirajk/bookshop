package order

import (
	"time"

	"context"

	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	next   Service
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingService{
			logger: logger,
			next:   next,
		}
	}
}

func (s loggingService) PlaceOrder(ctx context.Context, bookID string) (order Order, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "place_order",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.PlaceOrder(ctx, bookID)
}

func (s loggingService) GetUserOrders(ctx context.Context, userID string) (orders []Order, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "get_user_orders",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.GetUserOrders(ctx, userID)
}

func (s loggingService) CancelOrder(ctx context.Context, userID string, orderID string) (err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "cancel_order",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.CancelOrder(ctx, userID, orderID)
}
