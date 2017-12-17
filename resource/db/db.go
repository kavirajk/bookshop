package db

import (
	"errors"

	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/resource/config"
)

var (
	ErrMigrationsPathInvalid = errors.New("db: invalid migrations path")
)

type db struct {
	logger log.Logger
	config config.Datastore
	DB     *gorm.DB
}

func New(cfg config.Datastore, logger log.Logger) (*db, error) {
	db, err := gorm.Open(cfg.Driver, cfg.Address)
	if err != nil {
		return errors.Wrap(err, "db.New")
	}

	return &db{
		logger: logger,
		config: cfg,
		DB:     db,
	}
	logger.Log("event", "create db resource", "status", "success")
}

func (d *db) MigrateUp() error {

}

func (d *db) MigrateDown() error {

}
