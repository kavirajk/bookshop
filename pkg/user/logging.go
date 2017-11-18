package user

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

func (s loggingService) Register(ctx context.Context, nuser NewUser) (user User, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "register",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.Register(ctx, nuser)
}

func (s loggingService) Login(ctx context.Context, email, password string) (user User, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "login",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.Login(ctx, email, password)
}

func (s loggingService) AuthToken(ctx context.Context, token string) (user User, err error) {
	defer func(begin time.Time) {
		_ = s.logger.Log(
			"method", "login",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.AuthToken(ctx, token)
}

func (s loggingService) ResetPassword(ctx context.Context, key, newpass string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "reset-password",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())
	return s.next.ResetPassword(ctx, key, newpass)
}

func (s loggingService) ChangePassword(ctx context.Context, userID string, oldpass, newpass string) (err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "change-password",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return s.next.ChangePassword(ctx, userID, oldpass, newpass)
}

func (s loggingService) List(ctx context.Context, order string, limit, offset int) (users []User, total int, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "list",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return s.next.List(ctx, order, limit, offset)
}
