package db

import (
	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/resource/config"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

var (
	ErrMigrationsPathInvalid = errors.New("db: invalid migrations path")
)

type resource struct {
	logger log.Logger
	config config.Datastore
	DB     *gorm.DB
}

func New(cfg config.Datastore, logger log.Logger) (*resource, error) {
	db, err := gorm.Open(cfg.Driver, cfg.Address)
	if err != nil {
		return nil, errors.Wrap(err, "db.New")
	}

	logger.Log("event", "create db resource", "status", "success")

	return &resource{
		logger: logger,
		config: cfg,
		DB:     db,
	}, nil
}

func (r *resource) MigrateUp() error {
	// TODO(kaviraj): Replace with schema migrations
	// r.DB.AutoMigrate(&user.User{})
	return nil
}

func (r *resource) MigrateDown() error {
	return nil
}
