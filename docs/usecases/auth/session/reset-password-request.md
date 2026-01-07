# Reset Password Request

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.reset-password-request

## Summary

Initiates the password reset flow for users who have forgotten their password. Sends a secure reset link to the user's registered email address.

## Actor

- **Primary:** Anonymous User
- **Authorization:** None (public endpoint)

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| email | string | yes | valid email format |
| captcha_token | string | no | if CAPTCHA is enabled |

## Output

| Field | Type | Description |
|-------|------|-------------|
| success | bool | Always true (even if email not found) |
| message | string | Generic success message |

## Flow

### Validate

1. Validate email format
2. Verify CAPTCHA if enabled
3. Check rate limiting for this email/IP

### Execute

1. Look up user by email (via Person record)
2. If user not found, still return success (prevent enumeration)
3. If user found and active:
   - Generate secure random reset token (32 bytes, URL-safe base64)
   - Store token hash with expiry (1 hour) and user ID
   - Invalidate any existing reset tokens for this user
   - Send password reset email with link
4. Create audit log entry (whether found or not)

### Side Effects

- [ ] Send email: Password reset link (if user found)
- [ ] Create/update reset token record
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid email format |
| ErrCaptchaFailed | 400 | CAPTCHA verification failed |
| ErrTooManyRequests | 429 | Rate limit exceeded |

## Business Rules

1. Response must not reveal whether email exists (prevent enumeration)
2. Reset token expires after 1 hour
3. Only one active reset token per user (new request invalidates old)
4. Rate limit: max 3 requests per email per hour
5. Token should be single-use
6. Email should include:
   - Reset link
   - Expiry time
   - Warning if user didn't request this
   - IP address of requester

## Notes

- Use constant-time comparison when checking tokens
- Store only token hash in database, not the token itself
- Include request metadata (IP, user agent) in audit log
- Consider implementing notification if reset was requested but not completed
- Reset link should use HTTPS only
- Consider adding security questions as additional verification
