# Set Role Permission

Assigns permissions to a role. Can add or remove permissions.

> **type**: user_action

> **operation-id**: `set-role-permission`

> **access**: POST /auth/set-role-permission

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "role_id": 123,              // required, int64
    "permissions": ["string"]    // required, array of permission strings
}
```

## Output

```json
{
    "role_id": 123,
    "permissions": ["users:read", "users:write"]
}
```

## Execute

- Validate role_id exists
- Validate permission strings format
- Delete existing role permissions
- Insert new role permissions
- Return updated role permissions

## Error Scenarios

- `ROLE_NOT_FOUND`: Role does not exist
- `INVALID_PERMISSION_FORMAT`: Permission string format is invalid
