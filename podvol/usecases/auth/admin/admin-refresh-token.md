# Admin Refresh Token

Refreshes an expired access token using a valid refresh token.

> **type**: user_action

> **operation-id**: `admin-refresh-token`

> **access**: POST /auth/admin-refresh-token

> **actor**: admin

> **permissions**: none (authenticated admin)

## Input

```json
{
    "refresh_token": "string" // required
}
```

## Output

```json
{
    "access_token": "string",
    "access_token_expires_at": "2024-01-01T01:00:00Z",
    "refresh_token": "string",
    "refresh_token_expires_at": "2024-01-08T00:00:00Z"
}
```

## Execute

- Validate refresh token format
- Find session by refresh token
- Check if refresh token is not expired
- Check if session's admin is still active
- Generate new access token
- Generate new refresh token (rotate)
- Update session with new tokens
- Update session's last_used_at
- Return new tokens

## Error Scenarios

- `INVALID_REFRESH_TOKEN`: Refresh token not found or malformed
- `REFRESH_TOKEN_EXPIRED`: Refresh token has expired
- `ADMIN_DISABLED`: Associated admin account is disabled
