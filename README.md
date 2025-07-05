# DDD-Golang Todo Application

This project is a sample implementation of a Todo application using **Domain-Driven Design (DDD)** principles in Go. It demonstrates a clean architecture with clear separation of concerns, encapsulation, and best practices for structuring Go applications.

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
  common/             # Common/shared code (error helpers, etc.)
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
- **Common Layer**: Shared helpers and error utilities.

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
- Example error response:

```json
{
  "status_code": 400,
  "http_status": 400,
  "error_message": "Invalid title",
  "internal_reason": "Title validation failed",
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
  "status-code": 400,
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