```mermaid
---
title: RBAC Module - Entity Relationship Diagram
---
erDiagram
    %% ===========================================
    %% ROLE-BASED PATH (indirect permissions)
    %% Actor -> ActorRole -> Role -> RolePermission
    %% ===========================================

    roles {
        bigint id PK
        actor_type actor_type "user|admin|service_acc"
        varchar name UK "unique role name"
        timestamp created_at
        timestamp updated_at
    }

    role_permissions {
        bigint id PK
        bigint role_id FK
        varchar permission "e.g. users:read"
        timestamp created_at
    }

    actor_roles {
        bigint id PK
        actor_type actor_type
        varchar actor_id "polymorphic reference"
        bigint role_id FK
        timestamp created_at
    }

    %% ===========================================
    %% DIRECT PATH (bypasses roles)
    %% Actor -> ActorPermission
    %% ===========================================

    actor_permissions {
        bigint id PK
        actor_type actor_type
        varchar actor_id "polymorphic reference"
        varchar permission
        timestamp created_at
    }

    %% ===========================================
    %% RELATIONSHIPS
    %% ===========================================

    roles ||--o{ role_permissions : "grants"
    roles ||--o{ actor_roles : "assigned via"
```
