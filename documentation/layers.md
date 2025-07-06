## Layer Responsibilities

- **Domain Layer**: Contains pure business logic, aggregates, value objects, and domain services. No dependencies on other layers.
- **Application Layer**: The application layer, orchestrating use cases and coordinating domain logic and infrastructure. See [Use Case Explanation](#use-case-explanation) for details.

  **The Application Layer: The Orchestrator of Your System**

  The application layer is like the "orchestrator" or "conductor" of your system. It tells the domain layer what to do, coordinates business processes, and connects your core logic to the outside world (like web APIs or databases). The application layer doesn't contain business rules itself—it just makes sure the right things happen in the right order.

  **Detailed Explanation:**
  The application layer is responsible for defining and executing the use cases of your system. A use case is a specific business action or process that the system must support (e.g., "Create Todo", "List Todos", "Complete Todo"). This layer acts as the "glue" between the domain layer (which contains pure business logic) and the outside world (adapters and infrastructure).

  **Key Responsibilities:**
  - **Orchestrate Use Cases:** Each use case is implemented as a function or method (often in a `UseCase` struct) that coordinates the steps required to fulfill a business process. This may involve validating input, invoking domain logic, and interacting with repositories or external services. See [Use Case Explanation](#use-case-explanation) for more details.
  - **Coordinate Domain Logic:** The application layer calls domain models and domain services to enforce business rules. It does not contain business rules itself, but ensures they are applied in the correct order and context.
  - **Coordinate Infrastructure:** When a use case requires data persistence or external communication, the application layer interacts with infrastructure through well-defined interfaces (ports). This allows the application logic to remain decoupled from technical details.
  - **Input/Output Mapping:** The application layer receives input (often as command/query objects) from adapters, processes it, and returns output (often as response models) back to adapters for presentation.

  **Why is this important?**
  - **Separation of Concerns:** By keeping orchestration and coordination in the application layer, you prevent business logic from leaking into adapters or infrastructure, and vice versa.
  - **Testability:** Use cases can be tested in isolation by mocking domain and infrastructure dependencies.
  - **Flexibility:** Changes to how a use case is orchestrated (e.g., adding validation, changing the order of operations) can be made without touching domain or infrastructure code.
  - **Decoupling:** The application layer depends only on interfaces (ports), not concrete implementations, making it easy to swap out infrastructure or change delivery mechanisms (e.g., from HTTP to CLI).

  **Example in this project:**
  - The `TodoUseCase` struct in `application/usecase/` implements methods like `CreateTodoUseCase`, `ListTodosUseCase`, etc.
  - These methods receive command/query objects, validate them (possibly using domain services), call domain models to perform business logic, and interact with repositories via port interfaces.
  - The results are returned as application models (response objects) to the adapters.

> **The Adapters Layer: The Translator Between Worlds**
>
> The adapters layer is like the "translator" or "gateway" for your application. It lets the outside world (like users, web browsers, or command-line tools) talk to your system. Adapters convert incoming requests into a form your application understands, and turn responses from your core logic back into something the outside world can use. This keeps your core logic clean and focused on business rules, not technical details.

- **Adapters Layer**: Implements inbound interfaces (HTTP, CLI) and handles request/response mapping.

> **The Infrastructure Layer: The Toolbox and Plumbing**
>
> The infrastructure layer is like the "toolbox" or "plumbing" of your application. It handles all the technical details—like saving data to a database, sending emails, or talking to other services—so your core logic doesn't have to worry about them. By keeping these details separate, you can easily swap out tools or technologies without changing your business rules.

- **Infrastructure Layer**: Implements outbound interfaces (repositories, DB, etc.).