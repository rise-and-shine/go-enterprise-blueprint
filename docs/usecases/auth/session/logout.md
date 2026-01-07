# Logout

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.logout

## Summary

Terminates the current user session by invalidating the access and refresh tokens. Optionally can terminate all sessions for the user.

## Actor

- **Primary:** Authenticated User
- **Authorization:** Valid session required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| refresh_token | string | yes | valid refresh token |
| logout_all | bool | no | if true, invalidate all sessions |

## Output

| Field | Type | Description |
|-------|------|-------------|
| success | bool | Always true on success |
| sessions_terminated | int | Number of sessions terminated |
| message | string | Confirmation message |

## Flow

### Validate

1. Validate refresh token is present
2. Verify refresh token exists and is not already revoked
3. Extract user ID from token

### Execute

1. Load session associated with refresh token
2. If logout_all is true:
   - Find all active sessions for the user
   - Revoke all refresh tokens
   - Mark all sessions as terminated
3. If logout_all is false:
   - Revoke only the current refresh token
   - Mark current session as terminated
4. Add access token to blacklist (if using JWT blacklist)
5. Create audit log entry

### Side Effects

- [ ] Emit event: `UserLoggedOut`
- [ ] Invalidate refresh token(s)
- [ ] Add access token to blacklist
- [ ] Update session record(s)
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Missing refresh token |
| ErrUnauthorized | 401 | Invalid or expired refresh token |
| ErrSessionNotFound | 404 | Session already terminated |

## Business Rules

1. Logout should always succeed if token was ever valid (idempotent)
2. Access token should be immediately invalidated
3. Refresh token must be revoked to prevent reuse
4. If logout_all, all devices should be logged out
5. Client should discard tokens after logout

## Notes

- Implement token blacklist with TTL matching token expiry
- Consider using Redis for fast token blacklist lookups
- Logout should be graceful - don't fail if session is already gone
- Log logout events for security monitoring
- Consider notifying user of logout from other devices
