# Auth Module ERD

This diagram shows the Entity Relationship Diagram for the `auth` schema.

```mermaid
erDiagram
    admins {
        UUID id PK
        VARCHAR username UK
        VARCHAR password_hash
        BOOLEAN is_active
        TIMESTAMPTZ last_active_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    roles {
        BIGSERIAL id PK
        VARCHAR actor_type
        VARCHAR name UK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    role_permissions {
        BIGSERIAL id PK
        BIGINT role_id FK
        VARCHAR permission
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    actor_roles {
        BIGSERIAL id PK
        VARCHAR actor_type
        UUID actor_id
        BIGINT role_id FK
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    actor_permissions {
        BIGSERIAL id PK
        VARCHAR actor_type
        UUID actor_id
        VARCHAR permission
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    sessions {
        BIGSERIAL id PK
        VARCHAR actor_type
        UUID actor_id
        VARCHAR access_token
        TIMESTAMPTZ access_token_expires_at
        VARCHAR refresh_token
        TIMESTAMPTZ refresh_token_expires_at
        VARCHAR ip_address
        VARCHAR user_agent
        TIMESTAMPTZ last_used_at
        TIMESTAMPTZ created_at
        TIMESTAMPTZ updated_at
    }

    roles ||--o{ role_permissions : "has"
    roles ||--o{ actor_roles : "assigned via"
    admins ||--o{ actor_roles : "has (polymorphic)"
    admins ||--o{ actor_permissions : "has (polymorphic)"
    admins ||--o{ sessions : "has (polymorphic)"
```

## Table Descriptions

| Table | Description |
|-------|-------------|
| `admins` | Admin user accounts with authentication credentials |
| `roles` | Role definitions scoped by `actor_type` (e.g., admin, user) |
| `role_permissions` | Permissions granted to a role |
| `actor_roles` | Links actors (polymorphic via `actor_type` + `actor_id`) to roles |
| `actor_permissions` | Direct permissions granted to actors (bypassing roles) |
| `sessions` | Active sessions with access/refresh tokens for actors |

## Polymorphic Pattern

The auth module uses a polymorphic pattern for actors:
- `actor_type`: Identifies the type of actor (e.g., `"admin"`)
- `actor_id`: UUID reference to the actor's record in their respective table

This allows the same permission system to be used across different actor types (admins, future user types, etc.).
