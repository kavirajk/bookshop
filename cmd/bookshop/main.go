package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"

	"github.com/kavirajk/bookshop/resource/config"
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

	logger.Log("event", fmt.Sprintf("listening on: %s", cfg.Server.Address))
	logger.Log("error", (http.ListenAndServe(cfg.Server.Address, nil)))
}
