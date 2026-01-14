# Admin Login

Authenticates an admin user and creates a session with access and refresh tokens.

> **type**: user_action

> **operation-id**: `admin-login`

> **access**: POST /auth/admin-login

> **actor**: admin (unauthenticated)

> **permissions**: none (public endpoint)

## Input

```json
{
    "username": "string", // required, 3-50 chars
    "password": "string"  // required, min 8 chars
}
```

## Output

```json
{
    "admin": {
        "id": "string",
        "username": "string",
        "is_superadmin": true,
        "is_active": true,
        "last_active_at": "2024-01-01T00:00:00Z"
    },
    "session": {
        "access_token": "string",
        "access_token_expires_at": "2024-01-01T01:00:00Z",
        "refresh_token": "string",
        "refresh_token_expires_at": "2024-01-08T00:00:00Z"
    }
}
```

## Execute

- Validate input format
- Find admin by username
- Check if admin is active
- Verify password hash
- Generate access token (JWT, 1 hour expiry)
- Generate refresh token (random, 7 days expiry)
- Create session record with IP address and user agent
- Update admin's last_active_at timestamp
- Return admin info and session tokens

## Error Scenarios

- `INVALID_CREDENTIALS`: Username or password is incorrect
- `ADMIN_DISABLED`: Admin account is disabled
