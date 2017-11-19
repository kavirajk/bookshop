DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
default: build

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "build/bookshop" cmd/bookshop/main.go

.PHONE: test
test:
	@ TEST_PATH="${DIR}" go test -v -cover ./...
