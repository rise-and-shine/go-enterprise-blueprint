# Update Role

Updates an existing role's name.

> **type**: user_action

> **operation-id**: `update-role`

> **access**: POST /auth/update-role

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "id": 123,       // required, int64
    "name": "string" // required, 2-100 chars, unique
}
```

## Output

```json
{
    "id": 123,
    "name": "string",
    "updated_at": "2024-01-01T00:00:00Z"
}
```

## Execute

- Validate role ID and name
- Find role by ID
- Check if new name is unique (excluding current role)
- Update role name
- Return updated role

## Error Scenarios

- `ROLE_NOT_FOUND`: Role does not exist
- `ROLE_NAME_EXISTS`: Another role with this name already exists
- `INVALID_ROLE_NAME`: Role name does not meet requirements
