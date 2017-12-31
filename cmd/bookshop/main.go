package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"

	"github.com/kavirajk/bookshop/pkg/auth"
	"github.com/kavirajk/bookshop/resource/config"
	"github.com/kavirajk/bookshop/resource/db"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC, "caller", log.DefaultCaller)

	cfgPath := flag.String(
		"cfgpath", "",
		"config path to load the config values",
	)
	flag.Parse()

	if *cfgPath == "" {
		logger.Log("error", "cfgpath is missing. --help for more info")
		os.Exit(1)
	}

	cfg, err := config.FromFile(*cfgPath)
	if err != nil {
		logger.Log("error", err)
		os.Exit(1)
	}

	db, err := db.New(cfg.Datastore, logger)
	if err != nil {
		logger.Log("error", err)
	}

	pubKey, err := auth.LoadPublicKey("config/certs/public.pem")
	if err != nil {
		level.Error(logger).Log("name", "loading public key", "error", err)
		os.Exit(1)
	}

	privKey, err := auth.LoadPrivateKey("config/certs/private.pem")
	if err != nil {
		level.Error(logger).Log("name", "loading private key", "error", err)
		os.Exit(1)
	}

	tokenSvc := auth.NewTokenService("bs", privKey, pubKey, 24*10*time.Hour, logger)

	authSvc := auth.NewService(logger, tokenSvc, db)

	logger.Log("event", fmt.Sprintf("listening on: %s", cfg.Server.Address))

	logger.Log("error", (http.ListenAndServe(cfg.Server.Address, auth.MakeHTTPHandler(authSvc, logger))))
}
