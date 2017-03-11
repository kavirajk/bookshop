package catalog

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

func (s loggingService) Search(ctx context.Context, query string) (books []Book, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "search",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.Search(ctx, query)
}

func (s loggingService) List(ctx context.Context, order string, limit, offset int) (books []Book, total int, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "list",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.List(ctx, order, limit, offset)
}

func (s loggingService) Get(ctx context.Context, ID string) (book Book, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "get",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.Get(ctx, ID)
}
