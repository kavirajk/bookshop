package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	stdprometheus "github.com/prometheus/client_golang/prometheus"

	"context"

	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/kavirajk/bookshop/catalog"
	"github.com/kavirajk/bookshop/db/postgres"
	"github.com/kavirajk/bookshop/order"
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

	urepo, err := postgres.NewUserRepo(*dbSource)
	if err != nil {
		log.Fatalf("error creating user repo: %v\n", err)
	}

	crepo, err := postgres.NewCatalogRepo(*dbSource)
	if err != nil {
		log.Fatalf("error creating user repo: %v\n", err)
	}

	orepo, err := postgres.NewOrderRepo(*dbSource)
	if err != nil {
		log.Fatalf("error creating user repo: %v\n", err)
	}

	fieldKeys := []string{"method", "error"}

	var us user.Service
	us = user.NewService(urepo)
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

	var cs catalog.Service
	cs = catalog.NewService(crepo)
	cs = catalog.LoggingMiddleware(kitlog.NewContext(logger).With("component", "catalog"))(cs)
	cs = catalog.InstrumentingMiddleware(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "catalog_service",
			Name:      "request_count",
			Help:      "Number of requests received",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "catalog_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds",
		}, fieldKeys),
	)(cs)

	var os order.Service
	os = order.NewService(orepo)
	os = order.LoggingMiddleware(kitlog.NewContext(logger).With("component", "order"))(os)
	os = order.InstrumentingMiddleware(
		kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "order_service",
			Name:      "request_count",
			Help:      "Number of requests received",
		}, fieldKeys),
		kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "order_service",
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds",
		}, fieldKeys),
	)(os)

	httpLogger := kitlog.NewContext(logger).With("component", "http")
	mux := http.NewServeMux()

	userHandler := user.MakeHTTPHandler(ctx, us, httpLogger)
	catalogHandler := catalog.MakeHTTPHandler(ctx, cs, httpLogger)
	orderHandler := order.MakeHTTPHandler(ctx, os, httpLogger)

	mux.Handle("/users/v1/", userHandler)
	mux.Handle("/catalog/v1/", catalogHandler)
	mux.Handle("/order/v1/", orderHandler)

	mux.Handle("/metrics", stdprometheus.Handler())
	http.Handle("/", mux)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
