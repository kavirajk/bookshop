package main

import (
	"log"
	"net/http"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/db/postgres"
	"github.com/kavirajk/bookshop/user"
	"golang.org/x/net/context"
)

func main() {
	var logger kitlog.Logger
	logger = kitlog.NewLogfmtLogger(os.Stderr)
	ctx := context.Background()
	repo, err := postgres.NewUserRepo("user=kaviraj password=kaviraj dbname=bookstore sslmode=disable")
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
