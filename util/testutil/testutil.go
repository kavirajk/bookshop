package testutil

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
)

const (
	// To avoid reset actual db during testing in local.
	redisTestDB = 3

	testDBName = "bs-test"
)

var (
	redisPool *redis.Pool

	testDB *gorm.DB
)

// SetupRedis prepare redis server available for testing.
func SetupRedis(t *testing.T) (*redis.Pool, error) {
	dbOpt := redis.DialDatabase(redisTestDB)
	redisPool = &redis.Pool{
		MaxIdle:     5,
		Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", ":6379", dbOpt) },
		IdleTimeout: 240 * time.Second,
	}
	return redisPool, nil
}

// TearDownRedis handles all the cleanup work for redis in testing.
func TearDownRedis(t *testing.T) {
	if redisPool == nil {
		return
	}
	defer redisPool.Close()

	// flush the db
	conn := redisPool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")
}

func NewDB() *gorm.DB {
	gorm.Open()
}
