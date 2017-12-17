// Package config contains every config that are needed by
// all the services of bookstore.
package config // import "github.com/kavirajk/bookshop/resource/config"
import (
	"io/ioutil"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

// Config holds different types of config.
type Config struct {
	Datastore Datastore
	Server    Server
	Redis     Redis
}

// Datastore holds config related to data source URL, driver name,
// etc.
type Datastore struct {
	Driver         string `yaml:"driver"`
	Address        string `yaml:"address"`
	MigrationsPath string `yaml:"migrationsPath"`
}

// Server holds server config.
type Server struct {
	// address:port as a string
	Address string `yaml:"address"`
}

type Redis struct {
	//Address:port as a string
	Address string `yaml:"address"`
	DBNum   int    `yaml:"dbnum"`
}

// FromFile creates a Config by loading values from `path`.
func FromFile(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "config.new.failed to read config file %s", path)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, errors.Wrap(err, "config.new.unmarshal failed")
	}
	return &cfg, nil
}
