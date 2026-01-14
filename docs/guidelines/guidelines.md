# Guidelines

## Table of Contents

1. [Code Style](#code-style)
2. [Application Layers](#application-layers)
3. [Error Handling](#error-handling)
4. [Migrations](#migrations)
5. [Validation](#validation)
6. [Testing](#testing)
7. [Observability](#observability)

---

# Code Style

- **Style guides:** Effective Go, Go Code Review Comments, Uber Go Style Guide

## Package Naming

- Lowercase, single-word, matches directory name
- **Exception:** underscore allowed for use case layer (`create_user`)

```go
// Good
package user
package auth
package create_user  // use case layer only

// Bad
package userUtils
package common_helpers
```

## File Naming

Use snake_case: `user_repository.go`, `create_order.go`

## Struct Design

- Exported fields first, then unexported
- Group related fields
- Don't repeat default tags (e.g., if `bun` defaults to lowercase, omit tag)

## Comments

- **Don't** write obvious comments (`// GetUser gets a user`)
- Start exported item comments with the item name
- Keep comments up-to-date with the code

## Formatting and Linting

- Use `make fmt` and `make lint` to format and lint
- Don't push unformatted code or code with linting errors

---

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
- Define `OperationID` constant at top of file (e.g., `create-superadmin`)
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

---

# Error Handling

Use `github.com/code19m/errx` package for all error handling.

**Conventions:**

- Always handle errors, never ignore them
- Use `errx.Wrap()` to preserve stack traces
- Use error codes for programmatic error handling
- Never use `panic` in library/business code

## Wrapping Errors

Always wrap errors when returning:

```go
if err := doSomething(); err != nil {
    return errx.Wrap(err)
}
```

## Error Codes

Define codes instead of sentinel errors:

```go
var (
    CodeNotFound     = "NOT_FOUND"
    CodeUnauthorized = "UNAUTHORIZED"
    CodeInvalidInput = "INVALID_INPUT"
)
```

## Checking Errors

Use `errx.IsCodeIn` for error-based logic:

```go
if errx.IsCodeIn(err, CodeNotFound) {
    // handle not found
}
```

## Error Types

- Error types are defined **only at use case layer**
- All downstream errors should return `errx.T_Internal` (default)
- Use case layer knows its caller (actor) and assigns appropriate types
- Use `errx.WrapWithTypeOnCodes` to change type based on specific codes (e.g., when error is related to user input)

## Error Details

- Error Details are used to add additional context to the error
- The most proper layer to add details is the most downstream layer or when errx.New is called

---

# Migrations

Uses **goose** package.

## Commands

```bash
make migrate-create   # Create new migration
make migrate-up       # Apply all pending migrations
make migrate-down     # Rollback last migration
```

## File Naming Convention

Prefix with module name, use snake_case:

```bash
# Good
auth_init_schema
auth_add_user_roles
platform_init_taskmill

# Bad
init_schema
add_user_roles
```

## Structure

- All migrations in single `./migrations` folder
- Do NOT separate into subfolders by module

## Execution

Migrations run automatically on application startup (including production).

## Environment Variables

Required for Makefile commands:

```bash
POSTGRES_HOST
POSTGRES_PORT
POSTGRES_USER
POSTGRES_PASSWORD
POSTGRES_DB
POSTGRES_SSL
```

## Queries order

1. `CREATE TABLE` -- First always declare table creations
2. `CREATE INDEX` -- Then after each table declare it's indexes
3. `ALTER TABLE` -- Then after all table + index declarations declare table modifications such as adding foreign keys

## Rollback

Rollback works only for one migration at a time.
Always declare rollbacks in proper reverse order.
Always declare rollbacks with idempotency (with IF EXISTS keyword).

## Column types

For timestamp columns ALWAYS use `timestamptz` type
For non-unique strings ALWAYS use NOT NULL

---

# Validation

This document describes validation rules for each backend layer: Controller, Use Case, PBLC, and Infrastructure.

## Controller Validations

The controller layer validates only input parameters, not business logic.
For `user_action` use case types, input validation is handled automatically via `validate` tags on request structs using the `validator/v10` package.

## Use Case Validations

The use case layer handles two types of validation:

1. **Input validation** for cases not covered by the controller layer (e.g., cross-field validations that the validator package cannot handle).
2. **Business logic validation**, including verifying object IDs through the repository layer, checking object-level permissions, etc.

## PBLC Validations

The PBLC layer contains self-contained business logic and must validate its input fields strictly. However, this layer does not define error types as described in the error-handling guidelines. Instead, it returns error codes that the use case layer can handle programmatically if needed.

**Note:** Avoid defining too many error codes. Consider carefully whether a new error code is truly necessary for callers before adding one.

## Infrastructure Validations

The infrastructure layer handles data-layer validations with the following rules:

- **Database repositories**: Do not validate input parameters here. Validation is the responsibility of the use case and PBLC layers.
- **External services** (e.g., HTTP clients): Validate inputs strictly to catch potential errors as early as possible, before making calls to third-party services.

---

# Testing

## Mocking

- **Table-driven tests** preferred for comprehensive coverage
- Use `mockery` v2 for generating mocks (don't write mocks by hand)

```bash
make TODO:...
```

## Fake Objects

For external dependencies (APIs, services), write test doubles/fakes for integration tests:

```go
type FakeUserService struct {
    users map[int64]*User
}

func (f *FakeUserService) GetUser(ctx context.Context, id int64) (*User, error) {
    return f.users[id], nil
}
```

### Integration Testing

- Write **fake objects** (test doubles) for external dependencies
- Fakes provide controlled behavior for integration tests
- Located in `tests/` directory

TODO... explain our testing approach in GIVEN -> WHEN -> THEN pattern (after implementing good one).

---

# Observability

## Logging

- ALWAYS use our custom logging package (`github.com/rise-and-shine/pkg/observability/logger`).
- DO NOT use other logger packages.
