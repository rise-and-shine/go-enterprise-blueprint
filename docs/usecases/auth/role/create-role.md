# Create Role

Creates a new role in the RBAC system.

> **type**: user_action

> **operation-id**: `create-role`

> **access**: POST /auth/create-role

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "name": "string" // required, 2-100 chars, unique
}
```

## Output

```json
{
    "id": 123,
    "name": "string",
    "created_at": "2024-01-01T00:00:00Z"
}
```

## Execute

- Validate role name (length, format)
- Check if role with same name already exists
- Create role record
- Return created role

## Error Scenarios

- `ROLE_NAME_EXISTS`: Role with this name already exists
- `INVALID_ROLE_NAME`: Role name does not meet requirements
