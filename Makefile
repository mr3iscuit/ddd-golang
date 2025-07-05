# Makefile for ddd-golang

.PHONY: build run test lint swagger

build:
	go build -o build/bin/ddd-golang main.go

run:
	go run main.go

test:
	go test ./... -v

lint:
	golangci-lint run || true

swagger:
	swagger generate spec -o ./docs/swagger.json --scan-models 