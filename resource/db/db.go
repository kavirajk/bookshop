package db

import (
	"database/sql"

	"github.com/go-kit/kit/log"
	"github.com/jinzhu/gorm"
	"github.com/kavirajk/bookshop/resource/config"
	"github.com/mattes/migrate"
	"github.com/mattes/migrate/database/postgres"
	_ "github.com/mattes/migrate/source/file"
	"github.com/pkg/errors"
)

// New setup new database resource. It opens the connection, check the
// connectivity, run the migration and return the instance of *gorm.DB.
func New(cfg config.Datastore, logger log.Logger) (*gorm.DB, error) {
	db, err := sql.Open(cfg.Driver, cfg.Address)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database")
	}

	// check connectivity.
	if err := db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to connect to db instance")
	}

	// Run migrations
	if err := doMigration(cfg.MigrationsPath, db); err != nil {
		return nil, errors.Wrap(err, "failed doing migrations")
	}

	gdb, err := gorm.Open(cfg.Driver, db)
	if err != nil {
		return nil, errors.Wrap(err, "error opening database for gorm")
	}
	logger.Log("event", "create db resource", "status", "success")
	return gdb, nil

}

func doMigration(source string, db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return errors.Wrap(err, "failed to get postgres migrate driver")
	}
	mdb, err := migrate.NewWithDatabaseInstance(source, "bs", driver)
	if err != nil {
		return errors.Wrap(err, "failed to get migrate instance")
	}
	return mdb.Up()
}
