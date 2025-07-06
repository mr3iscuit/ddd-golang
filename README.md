# DDD-Golang Todo Application

 **What is this project about?**

 This is just a template project that is built using two powerful software design ideas: **Domain-Driven Design (DDD)** and **Hexagonal Architecture** (also called Ports & Adapters).

 - **Domain-Driven Design (DDD)** is about focusing your code on the real business or problem you want to solve. It means your code uses the same language and ideas as the people who know the business best.
 - **Hexagonal Architecture** is a way of organizing your code so that the "core" (the important business logic) is at the center, and all the ways of interacting with it (like web APIs, databases, or command-line tools) are kept at the edges. This makes your code easier to test, change, and grow.

 If you're new to these ideas, don't worry! This README will guide you step by step, explaining what each part does and why it matters.

## Overview

- [How to use this template ?](documentation/how_to_use_this_template.md)
- [Build & Run](#build--run)
- [Project Structure](#project-structure)
- [Naming Conventions](#naming-conventions)
- [JSON Field Naming](#json-field-naming)
- [Error Handling](#error-handling)
- [Flow Example](#flow-example)
- [Testing](documentation/testing.md)
- [How to Add a New Feature (with DDD Reasoning and Naming Context)](#how-to-add-a-new-feature-with-ddd-reasoning-and-naming-context)
- [Domain Service and Domain Model Explanation](#domain-service-and-domain-model-explanation)

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