> **Note:** This document is derived from the `go-enterprise-blueprint` project and may not fully reflect the current architecture of your project. However, it should be sufficient to describe the core architectural decisions. If you make significant changes to your architecture, consider updating this document to reflect those changes.

# Codebase Architecture

This document describes the code architecture and organization of the project.

## Overview

The project follows a **modular monolith** architecture with clean architecture principles. Each module is self-contained with its own domain, use cases, and controllers, while sharing common infrastructure through the `pkg/` layer and cross-module communication via the Portal pattern.

## Project Structure

```
go-enterprise-blueprint/
├── cmd/                    # Application entry points
├── config/                 # Configuration management
│   └── modconf/            # Module-specific configurations
├── docs/                   # Documentation (source of truth)
│   ├── architecture/       # Architecture documentation
│   ├── flows/              # Business/technical flow documentation
│   ├── usecases/           # Use case specifications
│   ├── integrations/       # External integration documentation
│   └── gen/dbdocs/         # Generated database documentation
├── internal/               # Private application code
│   ├── app/                # Application bootstrap and lifecycle
│   ├── modules/            # Business modules
│   └── portal/             # Cross-module communication interfaces. Contract between modules
├── pkg/                    # Shared packages (project-specific)
├── migrations/             # Database migrations (goose)
├── scripts/                # Utility bash scripts
└── tests/                  # Integration/E2E tests
```

## Architectural Layers

### 1. Entry Points (`cmd/`)

The command layer uses Cobra CLI framework to define application entry points:

- **run-all-in-one**: Runs all services in a single process (development/simple deployments)
- **HTTP server commands**: Individual HTTP server runners per module
- **CLI commands**: Module-specific management commands (e.g., `create-superuser`)
- **Cron manager**: Background task scheduler

### 2. Configuration (`config/`)

Configuration follows environment-based loading with YAML files:

```
config/
├── config.go              # Root configuration struct
├── modconf/               # Module-specific configs
│   └── auth.go            # Auth module config
├── local.yaml             # Local environment
├── staging.yaml           # Staging environment
└── production.yaml        # Production environment
```

Configuration is loaded via `github.com/rise-and-shine/pkg/cfgloader` with validation support.

### 3. Application Layer (`internal/app/`)

Responsible for:

- Application bootstrap and dependency wiring
- Service lifecycle management
- Logger initialization
- Global error handling setup

### 4. Modules (`internal/modules/`)

Each module is a self-contained business capability with the following internal structure:

```
internal/modules/{module}/
├── module.go              # Module initialization
├── domain/                # Domain layer
│   ├── container.go       # Domain container (DI)
│   ├── {entity}/          # Domain entities
│   │   ├── entity.go      # Entity definition
│   │   └── repo.go        # Repository interface
│   └── ...
├── service/               # Service layer
│   └── container.go       # Service container (DI)
├── usecase/               # Use case layer
│   └── {domain}/
│       └── {usecase}/
│           └── usecase.go # Use case implementation
├── ctrl/                  # Controller layer
│   ├── http/              # HTTP handlers
│   ├── cli/               # CLI commands
│   └── consumer/          # Event consumers
└── pblc/                  # Packaged Business Logic Components (reusable logic)
```

#### Current Modules

| Module      | Description                                           |
| ----------- | ----------------------------------------------------- |
| `auth`      | Authentication, authorization, user management, RBAC  |
| `esign`     | Electronic signature integration                      |
| `filestore` | File storage management                               |
| `doc`       | Documentation serving                                 |
| `platform`  | Platform utilities (taskmill, sentinel, integrations) |

### 5. Portal Layer (`internal/portal/`)

The Portal pattern provides controlled cross-module communication:

```go
// Container holds every module's portal interface
type Container struct {
    auth      auth.Portal
    doc       doc.Portal
    esign     esign.Portal
    filestore filestore.Portal
}
```

**Key principles:**

- Modules communicate only through Portal interfaces
- No direct imports between modules
- Portals expose only necessary functionality
- Enables module isolation and testing

## Module Internal Architecture

### Domain Layer

Contains business entities, value objects, and repository interfaces.

**Entity conventions:**

- Embed `pg.BaseModel` for PostgreSQL entities (provides `created_at`, `updated_at`)
- Use pointer types for nullable fields
- Don't repeat default struct tags (e.g., if `bun` defaults to lowercase field names)

```go
// Entity definition
type User struct {
    pg.BaseModel

    ID           int64   `json:"id"`
    Username     string  `json:"username"`
    Email        *string `json:"email"`       // nullable field
    PasswordHash string  `json:"-"`
}

// Repository interface - uses repogen generics
type UserRepo interface {
    repogen.Repo[UserFilter, User]
}
```

**Repository pattern:**

- Use `repogen` package for generic CRUD operations
- Prefer general-purpose methods over specialized queries
- Repository works with `bun.IDB` (supports both DB and Tx)

### Use Case Layer

Implements business logic following the use case pattern.

**Conventions:**

- **Document-first approach**: Always document use case in `docs/usecases/` before implementing
- Use case code should reference and reflect its documentation
- One file per use case, package name uses snake_case (team exception)
- Define `OperationID` as constant at top of file
- Enforce input validation and business logic rules
- Transaction management via UnitOfWork pattern at this layer

```go
package create_user

const OperationID = "auth.create-user"

// Use case definition
type UseCase ucdef.UserAction[*CreateUserIn, *CreateUserOut]

// Implementation
type useCase struct {
    dc domain.Container   // Domain dependencies
    sc service.Container  // Service dependencies
    pc portal.Container   // Cross-module dependencies
}

func (uc useCase) OperationID() string {
    return OperationID
}

func (uc useCase) Execute(ctx context.Context, in *CreateUserIn) (*CreateUserOut, error) {
    // 1. Validate input
    // 2. Check authorization
    // 3. Execute business logic
    // 4. Return result
}
```

#### Use Case Types (from `pkg/ucdef`)

| Type              | Trigger           | Example           |
| ----------------- | ----------------- | ----------------- |
| `UserAction`      | HTTP/gRPC request | CreateUser, Login |
| `EventSubscriber` | Domain event      | SendWelcomeEmail  |
| `AsyncTask`       | Cron/scheduler    | DailyReport       |
| `ManualCommand`   | CLI               | CreateSuperuser   |

### PBLC Layer (Packaged Business Logic Components)

Reusable business logic components called **only** from use cases:

- Deduplicate logic shared across multiple use cases
- Strict input validation (components don't know their callers)
- Suitable for OOP patterns (Strategy, State) when needed
- Flexible internal design within the layer

### Controller Layer

Adapts external interfaces to use cases.

**Conventions:**

- **One-to-one relationship** with use cases (each use case has exactly one controller)
- Keep this layer thin and simple
- Use generic components (e.g., `forward.ToUseCase`) instead of manual handlers
- No business logic in controllers

**HTTP Controller:**

```go
func (s *Server) initRoutes() {
    app := s.core.GetApp()
    // Use forward.ToUseCase for automatic request/response handling
    app.Post("/users", forward.ToUseCase(create_user.New(dc, sc, pc)))
}
```

**CLI Controller:**

```go
// Cobra commands for module-specific operations (manual_command use cases)
```

**Consumer Controller:**

```go
// Event consumers for async processing (event_subscriber use cases)
```

## Dependency Injection

Each layer uses container structs for dependency management:

```go
// Domain container
type Container struct {
    userRepo    user.UserRepo
    rbacRepo    rbac.RoleRepo
}

// Service container
type Container struct {
    hasher      hasher.Service
    mailer      mailer.Service
}

// Portal container (cross-module)
type Container struct {
    auth      auth.Portal
    esign     esign.Portal
}
```

Benefits:

- Explicit dependencies
- Easy testing with mocks
- Clear ownership boundaries

## HTTP Server Architecture

Base server with standard middleware stack:

```go
middlewares := []server.Middleware{
    middleware.NewRecoveryMW(cfg.Debug),     // Panic recovery
    middleware.NewTracingMW(),               // Distributed tracing
    middleware.NewTimeoutMW(cfg.Timeout),    // Request timeout
    middleware.NewMetaInjectMW(),            // Request metadata
    middleware.NewAlertingMW(),              // Error alerting
    middleware.NewLoggerMW(cfg.Debug),       // Request logging
    middleware.NewErrorHandlerMW(cfg.Debug), // Error formatting
}
```

Built on Fiber v2 for high performance.

## Error Handling

Uses `github.com/code19m/errx` for structured error handling with **error codes** (not sentinel errors):

```go
// Define error codes (not sentinel errors)
const (
    CodeNotFound     = "NOT_FOUND"
    CodeUnauthorized = "UNAUTHORIZED"
    CodeInvalidInput = "INVALID_INPUT"
)

// Create errors with codes
err := errx.New("user not found", errx.WithCode(CodeNotFound))

// Wrap errors to preserve stack trace
if err != nil {
    return errx.Wrap(err)
}

// Check error codes
if errx.IsCodeIn(err, CodeNotFound) {
    // handle not found
}
```

**Conventions:**

- Always handle errors, never ignore them
- Use `errx.Wrap()` to preserve stack traces
- Use error codes for programmatic error handling
- Never use `panic` in library/business code

## Key External Dependencies

| Package                         | Purpose                                           |
| ------------------------------- | ------------------------------------------------- |
| `github.com/rise-and-shine/pkg` | Shared enterprise packages (ucdef, repogen, etc.) |
| `github.com/spf13/cobra`        | CLI framework                                     |
| `github.com/gofiber/fiber/v2`   | HTTP framework                                    |
| `github.com/uptrace/bun`        | SQL ORM                                           |
| `github.com/jackc/pgx/v5`       | PostgreSQL driver                                 |
| `github.com/code19m/errx`       | Structured error handling with codes              |
| `go.opentelemetry.io/otel`      | Observability (tracing, metrics)                  |
| `go-resty/resty/v2`             | HTTP client for external integrations             |

## Database

- **Primary database:** PostgreSQL
- **ORM:** Bun (lightweight, type-safe)
- **Migrations:** Goose (SQL-based)
- **Connection pooling:** pgx with puddle

### Migration Conventions

- All migrations in single `./migrations/` folder (no subfolders)
- **Prefix with module name**: `{module}_{description}.sql`
- Migrations run **automatically on application startup** (including production)
- No manual DevOps intervention required for migrations

```bash
# Naming examples
auth_init_schema.sql        # Good - module prefix
auth_add_user_roles.sql     # Good - module prefix
platform_init_taskmill.sql  # Good - module prefix
init_schema.sql             # Bad - missing module prefix
```

Migration commands:

```bash
make migrate-create  # Create new migration
make migrate-up      # Apply migrations
make migrate-down    # Rollback migration
```

## Code Quality

- **Linting:** golangci-lint with comprehensive ruleset (`.golangci.yml`)
- **Formatting:** golangci-lint fmt
- **Style guides:** Effective Go, Go Code Review Comments, Uber Go Style Guide

```bash
make lint    # Run linter
make fmt     # Format code
go test -race ./...  # Run tests with race detector
```

## Testing

### Unit Testing

- **Table-driven tests** preferred for comprehensive coverage
- Use `mockery` v2 for generating mocks (don't write mocks by hand)
- Mocking via interfaces

```go
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name     string
        input    CreateUserIn
        wantErr  bool
        errCode  string
    }{
        {"valid user", CreateUserIn{...}, false, ""},
        {"invalid email", CreateUserIn{...}, true, CodeInvalidInput},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test logic
        })
    }
}
```

### Integration Testing

- Write **fake objects** (test doubles) for external dependencies
- Fakes provide controlled behavior for integration tests
- Located in `tests/` directory

## Diagrams

<!-- TODO: Add architecture diagrams -->
<!-- Reference: ./diagrams/ -->

## Related Documentation

- [Infrastructure Architecture](./infrastructure.md)
- [Use Cases](../usecases/)
- [Business Flows](../flows/)
