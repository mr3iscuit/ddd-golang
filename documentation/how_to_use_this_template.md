# How to use this template

This is a **template project** for quickly starting new DDD/Hexagonal Architecture applications. Here's how to modify and customize it for your own projects.

This project follows Domain-Driven Design (DDD) and Clean Architecture principles. Adding a new feature means respecting boundaries, naming, and responsibilities. Here's a step-by-step guide, with the *why* behind every step, what to create, where, and when.

If you are not familiar with Domain Driven Development consider viewing these files

- [Domain Driven Development](./ddd.md)

## Quick Start 

### 1. Start with the Domain Layer (Business Logic)
**Why:** The domain layer is the heart of your application. It models real business concepts, rules, and behaviors. By starting here, you ensure your code reflects real-world needs, not just technical requirements.

- **What to create:**
  - **Domain Model**: If your feature introduces a new business concept or modifies an existing one, create or update a struct in `domain/model/` (e.g., `Todo`, `User`).
  - **Domain Service**: If business logic spans multiple models or is not naturally part of an entity, create a service in `domain/service/`.
  - **Domain Event**: If something significant happens (e.g., a todo is completed), create an event in `domain/event/`.
- **Naming:**
  - Use business-centric names (e.g., `Todo`, `User`, `MarkAsImportant`).
  - Methods should express intent: `MarkAsCompleted`, `ArchiveTodo`, `PromoteToAdmin`.

This keeps business rules isolated and testable, and prevents accidental coupling to infrastructure or application logic.

### 2. Define Application Layer Contracts (Commands, Use Cases, Ports)
**Why:** The application layer orchestrates domain logic and coordinates tasks. It exposes use cases as interfaces (ports), and uses command/query objects for input.

- **What to create:**
  - **Command/Query**: If your use case requires structured input, create a command in `application/command/` (e.g., `MarkTodoImportantCommand`) or a query in `application/query/`.
  - **Use Case Implementation**: Implement the business process in `application/usecase/` (e.g., add a method to `TodoUseCase`).
  - **Inbound Port**: If this is a new use case, add a method to the relevant port interface in `application/port/` (e.g., `TodoUseCasePort`).
  - **Application Model**: If you need to return structured data, create or update a model in `application/model/`.
- **Naming:**
  - Commands: `XxxCommand` (e.g., `MarkTodoImportantCommand`)
  - Queries: `XxxQuery` (e.g., `ListImportantTodosQuery`)
  - Ports: `XxxPort` (e.g., `TodoUseCasePort`)
  - Use cases: `XxxUseCase` (e.g., `TodoUseCase`)
  - Application models: `XxxResponse`, `XxxRequest`
  - Ports define what the application can do, not how it's done.
  - **Do not create outbound ports (e.g., repository interfaces) unless your use case/domain logic requires interaction with external systems (DB, APIs, etc.).**
  - **Why:** This prevents unnecessary abstractions and keeps the codebase simple and focused. Only introduce outbound ports when the domain/application logic cannot proceed without them.

### 3. Update Adapters (HTTP, CLI, etc.)
**Why:** Adapters translate between the outside world (HTTP, CLI) and your application. They should only map requests/responses and call use cases.

- **What to create:**
  - **Handler/Endpoint/Command**: Add or update handler methods in `adapters/http/` or `adapters/cli/` (e.g., `todo_http_adapter.go`).
  - **Request/Response Mapping**: Map input to command/query objects, call the use case via the port, and map output to response models.
- **Naming:**
  - HTTP: `XxxHTTPAdapter` (e.g., `TodoHTTPAdapter`)
  - CLI: `XxxCLIAdapter` (e.g., `TodoCLIAdapter`)

Keeps delivery logic separate from business/application logic.
Makes it easy to add new interfaces (e.g., gRPC, GraphQL) without changing core logic.

### 4. Update Infrastructure (Persistence, External Services) — *If Needed*
**Why:** Infrastructure implements the technical details (DB, APIs) behind port interfaces. It should never contain business logic.

- **What to create:**
  - **Outbound Port (Repository Interface):** Only create a new outbound port in `application/port/` if your use case/domain logic requires it (e.g., saving or retrieving data).
  - **Repository Implementation:** Implement the outbound port in `infrastructure/repository/` (e.g., `todo_repository_postgres.go`).
  - **DB Mapper/Record:** Map between domain models and DB records (see `mapper.go`).
- **Naming:**
  - Outbound port: `XxxRepositoryPort` (interface, in application/port/)
  - Implementation: `XxxRepositoryPostgres` (in infrastructure/repository/postgres/)

Keeps persistence details swappable and testable.
Prevents leaking DB logic into business/application code.
**Do not create outbound ports or infrastructure code unless the domain/application logic requires it.**
**Why:** Avoids overengineering and keeps the codebase lean.

### 5. Update Tests
**Why:** Each layer should have its own tests, focusing on its responsibility.

- **What to create:**
  - **Domain tests:** Test business rules in `domain/model/` and `domain/service/`.
  - **Use case tests:** Test application logic in `application/usecase/`.
  - **Adapter tests:** Test request/response mapping in `adapters/http/` and `adapters/cli/`.
  - **Integration tests:** Test the full flow in `integration_test.go`.

### 6. Example: Adding 'Mark Todo as Important'

Suppose you want to add a feature to mark a todo as important.

**Domain Layer:**
- Add a field `isImportant` to `Todo` in `domain/model/todo.go`.
- Add a method `MarkAsImportant()` to `Todo`.

**Application Layer:**
- Add `MarkTodoImportantCommand` in `application/command/`.
- Add `MarkTodoImportantUseCase` method to `TodoUseCase` in `application/usecase/`.
- Add method to `TodoUseCasePort` in `application/port/`.
- Only if you need to persist this, add a method to `TodoRepositoryPort` (outbound port) and implement it in infrastructure.

**Adapters:**
- Add an endpoint/command in `adapters/http/todo_http_adapter.go` and/or `adapters/cli/todo_cli_adapter.go`.
- Map request to `MarkTodoImportantCommand`, call use case, map response.

**Infrastructure:**
- Only if persistence is needed, update DB schema and repository to persist `is_important`.

**Tests:**
- Add/extend tests for each layer.

### Naming Strategy: The "Why"
- **Commands/Queries:** Clearly state intent and input for a use case. (e.g., `MarkTodoImportantCommand`)
- **UseCase:** Encapsulates a single application action. (e.g., `MarkTodoImportantUseCase`)
- **Inbound Port:** Interface for use cases, enables easy mocking and swapping. (e.g., `TodoUseCasePort`)
- **Outbound Port:** Only when needed, interface for infrastructure dependencies. (e.g., `TodoRepositoryPort`)
- **Adapter:** Connects external world to application, keeps boundaries clear. (e.g., `TodoHTTPAdapter`)
- **Domain Model/Service:** Expresses business concepts and rules, not technical details.
- **Application Model:** Used for request/response at the application boundary.

### General 
- **Separation of Concerns:** Each layer has a single responsibility, making code easier to test, maintain, and evolve.
- **Explicit Boundaries:** Prevents accidental coupling and makes the system more robust to change.
- **Ubiquitous Language:** Naming matches business terms, improving communication with stakeholders.
- **Testability:** Each layer can be tested in isolation.
- **Avoid Premature Abstraction:** Only introduce interfaces and infrastructure when the domain/application logic requires them. This keeps the codebase simple and focused.

---

By following these steps and naming conventions, you ensure your codebase remains clean, maintainable, and true to DDD principles. Every name and boundary exists to make the system easier to change, and extend.