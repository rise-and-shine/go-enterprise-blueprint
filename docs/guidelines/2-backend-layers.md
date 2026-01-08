# Application Layers

## Controllers

- One-to-one relationship with use cases
- A use case cannot be called from multiple controllers
- Keep simple — use `forward.ToUseCase` for HTTP controllers
- Don't write manual controllers if possible

## Use Cases

Each use case has a type: `user_action`, `event_subscriber`, `async_task`, `manual_command`

- **Document first** — don't code before documenting
- Reference documentation in comments
- Define `OperationID` constant at top of file (e.g., `create-superuser`)
- Validate input (if not validated in controller)
- One file per use case
- If logic duplicates across use cases, move to PBLC layer
- Separate use cases for different actors for user_action types

## PBLC (Packaged Business Logic Component)

Reusable business logic called only from use cases.

- Components don't know their callers — validate inputs strictly
- Design freely (OOP patterns like State, Strategy when needed)
- Use for deduplicating logic across multiple use cases

## Repository

- Interfaces defined in domain layer
- Use `repogen` for PostgreSQL repositories
- Use `resty v2` for HTTP clients
- Prefer generalization over specialization
- Minimize custom methods — use repogen's general-purpose methods where possible

## Transaction Management

Manage with **UnitOfWork pattern** at use case layer.

Repository layer works with `bun.IDB` and is not responsible for transaction management.
