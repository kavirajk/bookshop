package testutil

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/resource/config"
	"github.com/kavirajk/bookshop/resource/db"
	"github.com/kavirajk/bookshop/util/pathutil"
	"github.com/mattes/migrate/database/postgres"
	"github.com/stretchr/testify/require"
)

const (
	// To avoid reset actual db during testing in local.
	redisTestDB = 3

	testDBName = "bs_test"
)

var (
	redisPool *redis.Pool

	testDB *gorm.DB
)

// SetupRedis prepare redis server available for testing.
func SetupRedis(t *testing.T) (*redis.Pool, error) {
	t.Helper()

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
	t.Helper()

	if redisPool == nil {
		return
	}
	defer redisPool.Close()

	// flush the db
	conn := redisPool.Get()
	defer conn.Close()
	conn.Do("FLUSHALL")
}

// NewDB setups new test db. It also handles running the migrations.
func NewDB(t *testing.T, logger log.Logger) (*gorm.DB, error) {
	root, err := pathutil.ProjectRoot()
	if err != nil {
		return nil, err
	}
	config, err := config.FromFile(root + "/config/test.yml")
	if err != nil {
		return nil, err
	}
	d, err := db.New(config.Datastore, logger)
	if err != nil {
		return nil, err
	}
	return d.Debug(), nil
}

// FlushDB used to drops tables created during the tests.
// It accepts optional tableNames to drop particular tables.
// By default it drops all the tables in test db.
func FlushDB(t *testing.T, db *gorm.DB, tableNames ...string) {
	require.NotNil(t, db)
	defer db.Close()

	if len(tableNames) == 0 {
		flushAllTables(t, db)
	}
	for _, table := range tableNames {
		flushSingleTable(t, db, table)
	}
	// delete schema_migrations table in any case
	flushSingleTable(t, db, postgres.DefaultMigrationsTable)
}

func flushSingleTable(t *testing.T, db *gorm.DB, tableName string) {
	require.NotNil(t, db)
	require.NoError(t, db.Exec(`DROP TABLE IF EXISTS `+tableName+` CASCADE`).Error)
}

func flushAllTables(t *testing.T, db *gorm.DB) {
	require.NotNil(t, db)
	var tableNames = make([]string, 0)

	require.NoError(t, db.Raw(`SELECT table_name FROM information_schema.tables WHERE table_schema=(SELECT CURRENT_SCHEMA())`).Scan(&tableNames).Error)

	if len(tableNames) == 0 {
		return
	}

	for _, table := range tableNames {
		flushSingleTable(t, db, table)
	}
}
