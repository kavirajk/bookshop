package main

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/kavirajk/bookshop/db/inmem"
	"github.com/kavirajk/bookshop/user"
	"golang.org/x/net/context"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	ctx := context.Background()
	repo := inmem.NewUserRepo()

	var userService user.Service
	userService = user.NewService(repo)
	userService = user.LoggingMiddleware(logger)(userService)

	httpLogger := log.NewContext(logger).With("component", "http")
	mux := http.NewServeMux()

	h := user.MakeHTTPHandler(ctx, userService, httpLogger)
	mux.Handle("/users/v1/", h)
	http.Handle("/", mux)

	http.ListenAndServe(":8080", nil)
}
