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
│   └── modules/            # Module-specific configurations
├── docs/                   # Documentation (source of truth)
├── internal/               # Private application code
│   ├── app/                # Application bootstrap and lifecycle
│   ├── modules/            # Business modules
│   └── portal/             # Cross-module communication interfaces. Contract between modules
├── pkg/                    # Shared packages (project-specific)
├── migrations/             # Database migrations (goose)
├── scripts/                # Utility bash scripts
└── tests/                  # Integration and E2E tests
```

## Architectural Layers

### 1. Entry Points (`cmd/`)

The command layer uses Cobra CLI framework to define application entry points:

- **Single binary** multiple commands
- **Hybrid** development and production deployments
- **Module-specific commands**: CLI commands for each module
- **CLI commands**: Module-specific management commands (e.g., `create-superadmin`)
- **run-all-in-one:** Runs all services with single command (development/simple deployments)

```bash
# List all available commands
go run ./cmd help
```

### 2. Configuration (`config/`)

Configuration follows environment-based loading with YAML files:

```
config/
├── config.go              # Root configuration struct
├── modules/               # Module-specific configs
│   ├── auth.go            # Auth module config
│   └── ...
├── production.yaml        # Production environment
├── staging.yaml           # Staging environment
├── dev.yaml               # Development environment
├── local.yaml.example     # Local environment example
├── local.yaml             # Local environment (gitignored)
└── test.yaml              # Test environment
```

Configuration is loaded via `github.com/rise-and-shine/pkg/cfgloader` with validation support.

### 3. Documentation (`docs/`)

The `docs/` folder serves as the source of truth for the project:

```
docs/
├── architecture/       # Architecture documentation (this file)
├── flows/              # Business/technical flow documentation
├── usecases/           # Use case specifications (API specs)
├── guidelines/         # Development guidelines and conventions
├── integrations/       # External integration documentation
├── uml/                # UML diagrams referenced from markdown files
├── misc/               # Miscellaneous files (PDFs, references)
└── gen/dbdocs/         # ? Generated database documentation (tbls)
```

We follow a **document-first approach**: documentation serves as the specification, and implementation follows the documentation.

### 4. Application Layer (`internal/app/`)

Responsible for:

- Application bootstrap and dependency wiring
- Service lifecycle management
- Graceful shutdown

### 5. Modules (`internal/modules/`)

Each module is a self-contained business capability with the following internal structure:

```
internal/modules/{module}/
├── module.go                  # Module initialization
├── domain/                    # Domain layer
│   ├── container.go           # Domain container (DI)
│   ├── {domain}/              # Domain entities
│   │   ├── entity.go          # Entity definition
│   │   └── repo.go            # Repository interface
│   └── ...
├── infra/                     # Infrastructure layer
│   └── postgres/
│       └── {domain}.go        # Implementation of domain Repo interfaces
│   └── http/
│       └── {client_name}.go   # Implementation of domain Client interfaces
│   └── ...
├── usecase/                   # Use case layer
│   └── {domain}/
│       └── {usecase}/
│           └── usecase.go     # Use case implementation
├── ctrl/                      # Controller layer
│   ├── http/                  # HTTP handlers
│   ├── cli/                   # CLI commands
│   └── consumer/              # Event consumers
├── pblc/                      # Packaged Business Logic Components (reusable logic)
└── portal/                    # Portal contract implementation (cross-module API)
```

### 6. Portal Layer (`internal/portal/`)

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
- Modules don't use other modules data directly (by joining).
- Modules don't know that we are using a single database.
- No transaction sharing between modules.

### 7. Testing (`tests/`)

TODO.... explain integration tests file structure

#### Current Modules

| Module     | Description                                           |
| ---------- | ----------------------------------------------------- |
| `auth`     | Authentication, authorization, user management, RBAC  |
| `esign`    | Electronic signature integration                      |
| `audit`    | Usecase based and Object based audit logging          |
| `platform` | Platform utilities (taskmill, sentinel, integrations) |
| ...        | ...                                                   |

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

    // Additional repository methods if needed
    // ...
}
```

**Repository pattern:**

- Use `repogen` package for generic CRUD operations
- Prefer general-purpose methods over specialized queries
- Repository works with `bun.IDB` (supports both DB and Tx)
- Search for client implementation from gitlab.mf.uz/go-lib/integration-sdk
- For http client implementations use resty v2. Example: TODO...

### Use Case Layer

Implements business logic following the use case pattern.

**Conventions:**

- **Document-first approach**: Always document use case in `docs/usecases/` before implementing
- One file per use case, package name uses snake_case (team exception)
- Define `OperationID` as constant at top of file

```go
package create_user

const OperationID = "create-user"

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
    // Validate input
    // Check preconditions
    // Execute business logic
    // ...
    // ...
    // Return result
}
```

#### Use Case Types (from `pkg/ucdef`)

| Type              | Trigger                                        | Example           |
| ----------------- | ---------------------------------------------- | ----------------- |
| `UserAction`      | HTTP/gRPC request                              | CreateUser, Login |
| `EventSubscriber` | Domain event                                   | SendWelcomeEmail  |
| `AsyncTask`       | Cron/scheduler or on-demand by other use cases | DailyReport       |
| `ManualCommand`   | CLI                                            | CreateSuperAdmin  |

For AsyncTask management we use taskmill framework `github.com/rise-and-shine/pkg/taskmill`

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
- Use generic components (e.g., `forward.ToUseCase`) instead of manual handlers where possible
- No business logic in controllers

## Dependency Injection

Each layer uses container structs (with getter methods) for dependency management:

```go
// Domain container
type Container struct {
    userRepo    user.UserRepo
    rbacRepo    rbac.RoleRepo
}

// PBLC container
type Container struct {
    hasher      hasher.Service
    mailer      mailer.Service
}

// Portal container (cross-module communications)
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

- Recovery
- Tracing
- Timeout
- Alerting
- Logging
- Error handling

Built on Fiber v2.
