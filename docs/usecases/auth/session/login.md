# Login

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.login

## Summary

Authenticates a user with their credentials and creates a new session. Returns access and refresh tokens for subsequent authenticated requests.

## Actor

- **Primary:** Anonymous User
- **Authorization:** None (public endpoint)

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| username | string | yes | min:3, max:50 |
| password | string | yes | min:1, max:128 |
| device_info | object | no | client metadata for session tracking |
| device_info.user_agent | string | no | max:500 |
| device_info.ip_address | string | no | valid IP format |
| device_info.device_id | string | no | max:100 |

## Output

| Field | Type | Description |
|-------|------|-------------|
| access_token | string | JWT access token for API authentication |
| refresh_token | string | Token for obtaining new access tokens |
| token_type | string | Always "Bearer" |
| expires_in | int | Access token expiry in seconds |
| refresh_expires_in | int | Refresh token expiry in seconds |
| user | object | Basic user information |
| user.id | int64 | User ID |
| user.username | string | Username |

## Flow

### Validate

1. Validate input format (username and password present)
2. Check rate limiting for the IP/username combination
3. Look up user by username

### Execute

1. Load User record by username
2. Verify user is active (is_active = true)
3. Verify password against stored hash using bcrypt
4. If authentication fails, increment failed attempt counter
5. If too many failures, temporarily lock the account
6. On success, reset failed attempt counter
7. Generate JWT access token with user claims
8. Generate refresh token and store in database
9. Create session record with device info
10. Update last_login_at timestamp
11. Create audit log entry

### Side Effects

- [ ] Emit event: `UserLoggedIn`
- [ ] Create session record
- [ ] Update last_login_at
- [ ] Create audit log
- [ ] Reset or increment failed login counter

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Missing or malformed credentials |
| ErrInvalidCredentials | 401 | Wrong username or password |
| ErrAccountLocked | 423 | Too many failed attempts |
| ErrAccountInactive | 403 | User account is deactivated |
| ErrTooManyRequests | 429 | Rate limit exceeded |

## Business Rules

1. Password comparison must be timing-safe to prevent timing attacks
2. Failed login attempts should be tracked per user and per IP
3. Account locks out after 5 consecutive failed attempts for 15 minutes
4. Access token expires in 15 minutes
5. Refresh token expires in 7 days
6. Session must be associated with device info for security tracking
7. Multiple concurrent sessions are allowed (configurable)

## Notes

- Implement brute-force protection with exponential backoff
- Consider adding CAPTCHA after multiple failed attempts
- Log all login attempts (success and failure) for security audit
- Support for MFA can be added as a second step (separate use case)
- Consider IP geolocation for anomaly detection
