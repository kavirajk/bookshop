package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"

	"context"
	kitlog "github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/db/postgres"
	"github.com/kavirajk/bookshop/user"
)

func main() {
	var (
		dbSource = flag.String("db", "", "Database source to connect to.")
	)
	flag.Parse()

	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	ctx := context.Background()

	if *dbSource == "" {
		fmt.Println("db argument is missing. Type --help for more info")
		os.Exit(1)
	}

	repo, err := postgres.NewUserRepo(*dbSource)
	if err != nil {
		log.Fatalf("error creating user repo: %v\n", err)
	}

	fieldKeys := []string{"method"}

	var us user.Service
	us = user.NewService(repo)
	us = user.LoggingMiddleware(kitlog.NewContext(logger).With("component", "user"))(us)
	us = user.InstrumentingMiddleware(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_count",
			Help:      "Number of requests received",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "user_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds",
		}, fieldKeys),
	)(us)

	httpLogger := kitlog.NewContext(logger).With("component", "http")
	mux := http.NewServeMux()

	h := user.MakeHTTPHandler(ctx, us, httpLogger)
	mux.Handle("/users/v1/", h)
	mux.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/", mux)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
