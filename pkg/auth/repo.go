package auth

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

var (
	DefaultTTL = time.Hour * 24 * 30 // 1 month
	// Used to mark the key(token)
	defaultValue = 1

	ErrTokenNotFound = errors.New("auth.repo: token not found")
	ErrTTLNotFound   = errors.New("auth.repo: ttl not found")
)

// Repo is the set of methods that manipulate store and retrieve of
// auth realted data in the persistant storage.
type Repo interface {
	// CheckToken checks if token present in redis store.
	CheckToken(rd *redis.Pool, token string) (bool, error)

	// GetTTL returns ttl of `token`. if token not
	// not present should return ErrTokenNotFound.
	GetTTL(rd *redis.Pool, token string) (time.Duration, error)

	// SaveToken put the token to redis store.
	SaveToken(rd *redis.Pool, token string, ttl time.Duration) error
}

// repo implements simple Repo interface.
type repo struct {
	logger log.Logger
}

// NewRepo create a simple Repo.
func NewRepo(logger log.Logger) Repo {
	return &repo{logger: logger}
}

// CheckToken checks whether given `token` is available in redis server.
func (r *repo) CheckToken(rd *redis.Pool, token string) (bool, error) {
	level.Debug(r.logger).Log("event", "CheckToken", "token", token)

	con := rd.Get()
	defer con.Close()

	exists, err := redis.Bool(con.Do("EXISTS", token))
	if err != nil {
		level.Error(r.logger).Log("event", "token_exists", "token", token, "error", err)
		return false, errors.Wrap(err, "auth.repo.checkToken.failed")
	}
	return exists, nil
}

func (r *repo) SaveToken(rd *redis.Pool, token string, ttl time.Duration) error {
	level.Debug(r.logger).Log("event", "SaveToken", "token", token)
	con := rd.Get()
	defer con.Close()

	if _, err := con.Do("SETEX", token, int(ttl.Seconds()), defaultValue); err != nil {
		level.Error(r.logger).Log("event", "SaveToken", "token", token, "err", err)
		return errors.Wrap(err, "auth.repo.SaveToken.failed")
	}
	return nil
}

func (r *repo) GetTTL(rd *redis.Pool, token string) (time.Duration, error) {
	level.Debug(r.logger).Log("event", "GetTTL", "token", token)
	con := rd.Get()
	defer con.Close()

	v, err := redis.Int(con.Do("TTL", token))
	if err != nil {
		level.Error(r.logger).Log("event", "GetTTL", "token", token, "err", err)
		return 0, errors.Wrap(err, "auth.repo.GetTTL.failed")
	}
	switch v {
	case -1:
		return 0, ErrTTLNotFound
	case -2:
		return 0, ErrTokenNotFound
	}
	return time.Duration(v) * time.Second, nil
}
