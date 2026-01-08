# Migrations

Uses **goose** package.

## Commands

```bash
make migrate-create   # Create new migration
make migrate-up       # Apply pending migrations
make migrate-down     # Rollback last migration
```

## Naming Convention

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
