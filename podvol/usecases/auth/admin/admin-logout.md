# Admin Logout

Logs out the current admin by invalidating their session.

> **type**: user_action

> **operation-id**: `admin-logout`

> **access**: POST /auth/admin-logout

> **actor**: admin

> **permissions**: none (authenticated admin)

## Input

No input required. Session is identified from the access token.

## Output

```json
{
    "success": true
}
```

## Execute

- Extract session ID from access token
- Delete the session record
- Return success

## Error Scenarios

- `SESSION_NOT_FOUND`: Session does not exist (already logged out)
