> **Note:** This document is derived from the `go-enterprise-blueprint` project and may not fully reflect the current architecture of your project. However, it should be sufficient to describe the core architectural decisions. If you make significant changes to your architecture, consider updating this document to reflect those changes.

# Infrastructure Architecture

This document describes the infrastructure components, deployment patterns, and operational aspects of the project.

## Overview

The application is designed as a **modular monolith** that can be deployed as:

- Single process (all-in-one) for simple deployments
- Multiple processes per module for scaled deployments

## Infrastructure Components

### Core Services

| Component     | Technology            | Purpose                  |
| ------------- | --------------------- | ------------------------ |
| Application   | Go 1.25+              | Business logic runtime   |
| Cache         | Redis                 | Session storage, caching |
| File Storage  | Minio                 | Media file storage       |
| Database      | PostgreSQL            | Primary data store       |
| Message Queue | PostgreSQL (taskmill) | Async event processing   |
| Pub/Sub       | PostgreSQL (planned)  | Pub Sub messaging        |

## Deployment Patterns

### Development (All-in-One)

Single process running all components:

```bash
go run ./cmd run-all-in-one
```

Suitable for:

- Local development
- Small deployments
- Testing

### Production & Staging (Distributed)

Separate processes per concern:

```bash
# Single runner for each module.
# All processes of a module run together, e.g. cron manager, event consumer, http server
./app run-auth
./app run-audit
./app run-platform

# Also if scaling requires processes of a module to run separately
./app run-auth-http
./app run-auth-taskmill-worker
./app run-auth-tasmilll-scheduler

# CLI commands (one-off)
./app auth create-superuser
```

### Container Deployment

- Multi-Stage build
- Single Docker image

## Configuration Management

### Environment-Based Configuration

```
config/
├── local.yaml        # Local development
├── staging.yaml      # Staging environment
├── production.yaml   # Production environment
└── ...
```

Configuration is loaded based on `ENVIRONMENT` environment variable:

- `ENVIRONMENT=local` loads `local.yaml`
- `ENVIRONMENT=staging` loads `staging.yaml`
- `ENVIRONMENT=production` loads `production.yaml`
- ...

### Secrets Management

Sensitive values are marked with `secret:"true"` tag and masked in logs

```go
type Config struct {
    SecretField string `yaml:"secret_field" secret:"true"`
}
```

### Configuration Validation

All configuration is validated at startup using `go-playground/validator/v10`:

```go
type Config struct {
    Logger logger.Config `yaml:"logger" validate:"required"`
}
```

## Observability Stack

### Logging

- **Library:** Custom wrapper around zap (structured logging)
- **Format:** JSON in production, pretty in local development
- **Levels:** debug, info, warn, error

Logger configuration:

```yaml
logger:
  level: info # debug | info | warn | error
  encoding: json # json | pretty
```

### Distributed Tracing

- **Standard:** OpenTelemetry

### Metrics

- **Standard:** OpenTelemetry Metrics

### Application Error Alerting

- **Library:** Sentinel
