# Makefile for ddd-golang

.PHONY: build run run-built test lint swagger clean docker-build docker-up docker-down docker-logs

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

# Docker commands
docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-restart:
	docker compose restart

docker-clean:
	docker compose down -v --remove-orphans

# Database migration commands
MIGRATE_CMD=migrate -path ./migrations -database "postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable"

migrate-up:
	$(MIGRATE_CMD) up

migrate-down:
	$(MIGRATE_CMD) down 1

migrate-force:
	$(MIGRATE_CMD) force

migrate-version:
	$(MIGRATE_CMD) version

# Integration test with Postgres (dev env)
integration-test:
	@echo "Checking Docker..."
	@docker info > /dev/null 2>&1 || (echo "Docker is not running or accessible. Please start Docker and try again." && exit 1)
	@echo "Starting Postgres service..."
	docker compose up -d postgres
	@echo "Waiting for Postgres to be ready..."
	@sleep 10
	@echo "Running integration tests..."
	@eval $(grep -v '^#' ./.env | xargs -0) DB_HOST=localhost DB_PORT=5432 DB_USER=todo_user DB_PASSWORD=todo_password DB_NAME=todo_db go test -v integration_test.go
	@echo "Running repository tests..."
	@eval $(grep -v '^#' ./.env | xargs -0) DB_HOST=localhost DB_PORT=5432 DB_USER=todo_user DB_PASSWORD=todo_password DB_NAME=todo_db go test -v ./infrastructure/repository/postgres
	@echo "Cleaning up..."
	docker compose down
