# Delete Role

Deletes a role and all its associated permissions and actor assignments (cascade).

> **type**: user_action

> **operation-id**: `delete-role`

> **access**: POST /auth/delete-role

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "id": 123 // required, int64
}
```

## Output

```json
{
    "success": true
}
```

## Execute

- Validate role ID
- Find role by ID
- Delete role (cascade deletes role_permissions and actor_roles via FK)
- Return success

## Error Scenarios

- `ROLE_NOT_FOUND`: Role does not exist
