# Disable Admin

Disables an admin account, preventing login. Also terminates all active sessions.

> **type**: user_action

> **operation-id**: `disable-admin`

> **access**: POST /auth/disable-admin

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "id": "uuid-string" // required
}
```

## Output

```json
{
    "id": "uuid-string",
    "username": "string",
    "is_active": false,
    "sessions_terminated": 3
}
```

## Execute

- Validate admin ID exists
- Check that admin is not the last active superadmin
- Set admin's `is_active` to false
- Delete all sessions for this admin
- Return result with count of terminated sessions

## Error Scenarios

- `ADMIN_NOT_FOUND`: Admin does not exist
- `ADMIN_ALREADY_DISABLED`: Admin is already disabled
- `CANNOT_DISABLE_LAST_SUPERADMIN`: Cannot disable the last active superadmin
