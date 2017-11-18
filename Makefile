default: build

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o "build/bookshop" cmd/bookshop/main.go
