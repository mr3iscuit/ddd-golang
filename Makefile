# Makefile for ddd-golang

.PHONY: build run run-built test lint swagger clean

build:
	go build -o build/bin/ddd-golang main.go

run:
	go run main.go

run-built: build
	./build/bin/ddd-golang

test:
	go test ./... -v

lint:
	golangci-lint run || true

swagger:
	swagger generate spec -o ./docs/swagger.json --scan-models

clean:
	rm -rf build/ 