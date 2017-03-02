package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/db/postgres"
	"github.com/kavirajk/bookshop/user"
	"golang.org/x/net/context"
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

	var userService user.Service
	userService = user.NewService(repo)
	userService = user.LoggingMiddleware(logger)(userService)

	httpLogger := kitlog.NewContext(logger).With("component", "http")
	mux := http.NewServeMux()

	h := user.MakeHTTPHandler(ctx, userService, httpLogger)
	mux.Handle("/users/v1/", h)
	http.Handle("/", mux)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
