# Get Role Permissions

Retrieves all permissions assigned to a specific role.

> **type**: user_action

> **operation-id**: `get-role-permissions`

> **access**: GET /auth/get-role-permissions

> **actor**: admin

> **permissions**: `superadmin`

## Input

Query parameters:

- `role_id`: int64, required

## Output

```json
{
    "role_id": 123,
    "role_name": "string",
    "permissions": ["users:read", "users:write", "docs:read"]
}
```

## Execute

- Validate role_id
- Find role by ID
- Query role permissions
- Return role with permissions list

## Error Scenarios

- `ROLE_NOT_FOUND`: Role does not exist
