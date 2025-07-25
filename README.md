# DDD-Golang Todo Application

This project is a sample implementation of a Todo application using **Domain-Driven Design (DDD)** principles in Go. It demonstrates a clean architecture with clear separation of concerns, encapsulation, and best practices for structuring Go applications.

## Build & Run

### Docker Setup (Recommended)

**Build and start the application with PostgreSQL:**
```sh
make docker-build
make docker-up
```

**View logs:**
```sh
make docker-logs
```

**Stop the application:**
```sh
make docker-down
```

**Restart the application:**
```sh
make docker-restart
```

**Clean up Docker resources:**
```sh
make docker-clean
```

The application will be available at `http://localhost:8080` and PostgreSQL at `localhost:5432`.

### Available Makefile Operations

**Build the project binary:**
```sh
make build
```
The binary will be created at `build/bin/ddd-golang`.

**Run the application (using go run):**
```sh
make run
```

**Build and run the binary:**
```sh
make run-built
```

**Run tests:**
```sh
make test
```

**Run linting:**
```sh
make lint
```

**Generate Swagger documentation:**
```sh
make swagger
```

**Clean build artifacts:**
```sh
make clean
```

### Manual Commands

To run the application directly:
```sh
go run main.go
```

Or run the built binary directly:
```sh
./build/bin/ddd-golang
```

## Project Structure

```
ddd-golang/
  adapters/           # Adapters for CLI and HTTP (inbound interfaces)
    cli/              # CLI adapter
    http/             # HTTP adapter (REST API)
  application/        # Application layer (use cases, commands, queries, ports, models)
    command/          # Command objects for use cases
    model/            # Application models (response/request models, error responses)
    port/             # Port interfaces (inbound/outbound)
    query/            # Query objects for use cases
    usecase/          # Use case implementations
  domain/             # Domain layer (entities, value objects, domain services)
    event/            # Domain events
    model/            # Domain models (pure business logic)
    service/          # Domain services
  infrastructure/     # Infrastructure layer (repositories, DB, etc.)
    repository/       # In-memory repository implementation
  main.go             # Application entry point
  integration_test.go # Integration tests
  README.md           # This file
```

## Layer Responsibilities

- **Domain Layer**: Contains pure business logic, aggregates, value objects, and domain services. No dependencies on other layers.
- **Application Layer**: Contains use cases, application models (formerly DTOs), commands, queries, and port interfaces. Coordinates domain logic and infrastructure.
- **Adapters Layer**: Implements inbound interfaces (HTTP, CLI) and handles request/response mapping.
- **Infrastructure Layer**: Implements outbound interfaces (repositories, DB, etc.).

## Naming Conventions

- **Port**: Interface for inbound/outbound communication (e.g., `TodoUseCasePort`, `TodoRepositoryPort`).
- **UseCase**: Application service implementing a use case (e.g., `TodoUseCase`).
- **Adapter**: Inbound adapter (e.g., `TodoHTTPAdapter`, `TodoCLIAdapter`).
- **Model**: Application-layer models for request/response and error responses (formerly DTOs).
- **DomainErrorPort**: Interface for domain errors with getter methods.

## JSON Field Naming

All JSON fields in API responses use **kebab-case** (e.g., `created-at`, `completed-at`, `error-message`).

## Error Handling

- All errors returned from use cases and adapters are domain errors implementing the `DomainErrorPort` interface.
- HTTP and CLI adapters map domain errors to structured error response models (see `application/model/error_response.go`).
- Errors are organized into logical groups with specific error code ranges:

### Error Code Ranges

- **1000-1999**: Validation errors (title, description, priority validation)
- **2000-2999**: Not found errors (todo not found)
- **3000-3999**: Operation errors (cannot complete/archive)
- **4000-4999**: Repository errors (database operations)
- **5000-5999**: HTTP errors (JSON parsing)
- **9000-9999**: Test errors (for testing purposes)

### Example Error Response

```json
{
  "error-code": 1001,
  "http-status": 400,
  "error-message": "Invalid title",
  "internal-reason": "Title validation failed",
  "details": {"max_length": "100"}
}
```

## Flow Example

1. **HTTP Request** → HTTP Adapter parses request and calls the use case via the port interface.
2. **Use Case** → Coordinates domain logic, calls domain models/services, and repositories.
3. **Domain Model** → Encapsulates business rules and validation.
4. **Repository** → Persists and retrieves domain models.
5. **Response** → Use case returns application model (response model) or domain error.
6. **HTTP Adapter** → Maps result to JSON response (kebab-case fields).

## Example: Create Todo (HTTP)

- **Request:**

```json
{
  "title": "Buy milk",
  "description": "Get 2 liters of milk",
  "priority": "high"
}
```

- **Response:**

```json
{
  "id": "..."
}
```

- **Error Response:**

```json
{
  "error-code": 1001,
  "http-status": 400,
  "error-message": "Invalid title",
  "internal-reason": "Title validation failed"
}
```

## Testing

- Run all tests:
  ```sh
  go test ./... -v
  ```
- Integration tests use a real HTTP server and test the full flow using `curl`.

## DDD Principles Applied

- **Encapsulation**: Domain models have private fields and public getters/setters.
- **Separation of Concerns**: Each layer has a clear responsibility.
- **Ports & Adapters**: All communication is via interfaces (ports), implemented by adapters.
- **No DTOs in Domain**: Domain layer only uses pure domain models; application models (formerly DTOs) are used in the application and adapter layers.

## Contributing

Feel free to open issues or PRs to improve the structure, add features, or suggest improvements! 