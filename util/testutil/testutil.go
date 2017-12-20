package testutil

import (
	"os"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/resource/config"
	"github.com/kavirajk/bookshop/resource/db"
	"github.com/kavirajk/bookshop/util/pathutil"
	_ "github.com/lib/pq"
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

func NewDB() (*gorm.DB, error) {
	root, err := pathutil.ProjectRoot()
	if err != nil {
		return nil, err
	}
	config, err := config.FromFile(root + "/config/test.yml")
	if err != nil {
		return nil, err
	}
	d, err := db.New(config.Datastore, log.NewLogfmtLogger(os.Stdout))
	if err != nil {
		return nil, err
	}
	return d.DB, nil
}
