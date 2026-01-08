> **Note:** This document is derived from the `go-enterprise-blueprint` project and may not fully reflect the current architecture of your project. However, it should be sufficient to describe the core architectural decisions. If you make significant changes to your architecture, consider updating this document to reflect those changes.

# Infrastructure Architecture

This document describes the infrastructure components, deployment patterns, and operational aspects of the project.

## Overview

The application is designed as a **modular monolith** that can be deployed as:

- Single process (all-in-one) for simple deployments
- Multiple processes per module for scaled deployments
- Containerized services for cloud-native environments

## Infrastructure Components

### Core Services

| Component     | Technology      | Purpose                  |
| ------------- | --------------- | ------------------------ |
| Application   | Go 1.25+        | Business logic runtime   |
| Database      | PostgreSQL      | Primary data store       |
| Cache         | Redis (planned) | Session storage, caching |
| Message Queue | (planned)       | Async event processing   |

### Supporting Services

| Component     | Technology    | Purpose                         |
| ------------- | ------------- | ------------------------------- |
| Reverse Proxy | Nginx/Traefik | Load balancing, TLS termination |
| Observability | OpenTelemetry | Distributed tracing, metrics    |
| Alerting      | Sentinel      | Error tracking and alerts       |
| Logging       | Zap           | Structured JSON logging         |

## Database Architecture

### PostgreSQL Configuration

- **Version:** PostgreSQL 15+
- **Connection pooling:** pgx with puddle
- **ORM:** Bun (lightweight SQL-first)

Connection string format:

```
postgres://{user}:{password}@{host}:{port}/{database}?sslmode={ssl}
```

Environment variables:

```bash
POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_HOST=
POSTGRES_PORT=5432
POSTGRES_DB=
POSTGRES_SSL=disable  # disable | require | verify-full
```

### Migration Strategy

- **Tool:** Goose
- **Migration table:** `_migrations`
- **Location:** `./migrations/`
- **Execution:** Automatic on application startup (including production)
- No manual DevOps intervention required

Commands (for manual execution):

```bash
make migrate-create   # Create new migration
make migrate-up       # Apply pending migrations
make migrate-down     # Rollback last migration
```

Migration naming convention:

```
{timestamp}_{module}_{description}.sql
# Example: 20251214105116_auth_init_schema.sql

# Good - prefixed with module name
auth_init_schema
auth_add_user_roles
platform_init_taskmill

# Bad - missing module prefix
init_schema
add_user_roles
```

### Database Documentation

Database documentation is generated using `tbls` and stored in:

```
docs/gen/dbdocs/
```

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

### Production (Distributed)

Separate processes per concern:

```bash
# HTTP API servers (can be scaled horizontally)
./app run-auth-http
./app run-platform-http

# Background workers
./app run-cron-manager
./app run-event-consumer

# CLI commands (one-off)
./app auth create-superuser
```

### Container Deployment

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/bin/app ./cmd

# Runtime stage
FROM alpine:latest
COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/config /app/config
ENTRYPOINT ["/app/app"]
```

## Configuration Management

### Environment-Based Configuration

```
config/
├── local.yaml        # Local development
├── staging.yaml      # Staging environment
└── production.yaml   # Production environment
```

Configuration is loaded based on `APP_ENV` environment variable:

- `APP_ENV=local` loads `local.yaml`
- `APP_ENV=staging` loads `staging.yaml`
- `APP_ENV=production` loads `production.yaml`

### Secrets Management

Sensitive values are marked with `secret:"true"` tag and should be:

- Loaded from environment variables
- Stored in secrets manager (Vault, AWS Secrets Manager)
- Never committed to version control

```go
type Config struct {
    SecretField string `yaml:"secret_field" secret:"true"`
}
```

### Configuration Validation

All configuration is validated at startup using `go-playground/validator`:

```go
type Config struct {
    Logger logger.Config `yaml:"logger" validate:"required"`
}
```

## Observability Stack

### Logging

- **Library:** Zap (structured logging)
- **Format:** JSON in production, console in development
- **Levels:** debug, info, warn, error

Logger configuration:

```yaml
logger:
  level: info # debug | info | warn | error
  format: json # json | console
  output: stdout # stdout | file path
```

### Distributed Tracing

- **Standard:** OpenTelemetry
- **Propagation:** W3C Trace Context
- **Exporters:** OTLP (configurable)

Tracing is injected via middleware:

```go
middleware.NewTracingMW()
```

### Metrics

- **Standard:** OpenTelemetry Metrics
- **Export:** Prometheus format

Key metrics:

- HTTP request duration
- HTTP request count by status
- Database query duration
- Background job execution

### Error Alerting

- **Library:** Sentinel
- **Features:** Error grouping, rate limiting, context capture

Alerting middleware:

```go
middleware.NewAlertingMW()
```

## Network Architecture

### HTTP Server Configuration

```yaml
http_server:
  host: 0.0.0.0
  port: 8080
  debug: false
  handle_timeout: 30s
```

### Middleware Stack

Request flow through middleware (in order):

1. **Recovery** - Panic recovery, prevents crashes
2. **Tracing** - Distributed tracing context
3. **Timeout** - Request timeout enforcement
4. **MetaInject** - Request metadata (request ID, etc.)
5. **Alerting** - Error capture for alerting
6. **Logger** - Request/response logging
7. **ErrorHandler** - Standardized error responses

### API Gateway Integration

For production deployments behind API gateway:

- Health check endpoint: `GET /health`
- Readiness endpoint: `GET /ready`
- Metrics endpoint: `GET /metrics`

## Security

### Authentication

- JWT-based authentication
- Token refresh mechanism
- Session management

### Authorization

- RBAC (Role-Based Access Control)
- Permission-based checks
- Actor types: `user`, `admin`, `service_acc`

### TLS/SSL

- TLS termination at load balancer/reverse proxy
- Internal communication can be plain HTTP
- Certificate management via cert-manager (Kubernetes) or similar

### Environment Security

```bash
# Never commit these files
.env
*.pem
*.key
credentials.json
```

## Scaling Considerations

### Horizontal Scaling

HTTP servers are stateless and can be scaled horizontally:

- Load balancer distributes requests
- Session state stored externally (Redis)
- Database connection pooling per instance

### Vertical Scaling

For single-instance deployments:

- Increase container resources
- Tune database connection pool size
- Adjust worker concurrency

### Database Scaling

- Read replicas for read-heavy workloads
- Connection pooling (PgBouncer) for high concurrency
- Table partitioning for large tables

## Backup and Recovery

### Database Backups

- Daily automated backups
- Point-in-time recovery enabled
- Backup retention policy: 30 days

### Disaster Recovery

- Multi-AZ deployment for high availability
- Database failover automation
- Regular disaster recovery testing

## Monitoring and Alerts

### Key Metrics to Monitor

| Metric               | Warning | Critical |
| -------------------- | ------- | -------- |
| HTTP 5xx rate        | > 1%    | > 5%     |
| Response latency P99 | > 1s    | > 5s     |
| Database connections | > 80%   | > 95%    |
| CPU usage            | > 70%   | > 90%    |
| Memory usage         | > 80%   | > 95%    |

### Health Checks

```go
// Liveness - is the process running?
GET /health -> 200 OK

// Readiness - can it accept traffic?
GET /ready -> 200 OK (checks DB connection)
```

## Infrastructure Diagrams

<!-- TODO: Add infrastructure diagrams -->

### Deployment Diagram

```
┌────────────────────────────────────────────────────────┐
│                      Load Balancer                     │
│                    (Nginx/Traefik)                     │
└─────────────────────────┬──────────────────────────────┘
                          │
          ┌───────────────┼───────────────┐
          │               │               │
          ▼               ▼               ▼
    ┌──────────┐   ┌──────────┐   ┌──────────┐
    │ App #1   │   │ App #2   │   │ App #3   │
    │ (HTTP)   │   │ (HTTP)   │   │ (HTTP)   │
    └────┬─────┘   └────┬─────┘   └────┬─────┘
         │              │              │
         └──────────────┼──────────────┘
                        │
                        ▼
              ┌──────────────────┐
              │   PostgreSQL     │
              │   (Primary)      │
              └────────┬─────────┘
                       │
                       ▼
              ┌──────────────────┐
              │   PostgreSQL     │
              │   (Replica)      │
              └──────────────────┘
```

## Related Documentation

- [Codebase Architecture](./codebase.md)
- [Database Schema](../gen/dbdocs/)
- [Integration Documentation](../integrations/)
