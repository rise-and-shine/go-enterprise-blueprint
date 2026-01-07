# Reset Password Confirm

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.reset-password-confirm

## Summary

Completes the password reset flow by validating the reset token and setting a new password. This is the second step after reset-password-request.

## Actor

- **Primary:** Anonymous User (with valid reset token)
- **Authorization:** Valid reset token required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| token | string | yes | valid reset token |
| new_password | string | yes | min:8, max:128, complexity requirements |
| confirm_password | string | yes | must match new_password |

## Output

| Field | Type | Description |
|-------|------|-------------|
| success | bool | True on successful password reset |
| message | string | Confirmation message |

## Flow

### Validate

1. Validate token is present and well-formed
2. Verify new_password matches confirm_password
3. Validate new password meets complexity requirements:
   - Minimum 8 characters
   - At least one uppercase letter
   - At least one lowercase letter
   - At least one digit
4. Hash the token and look up in database
5. Verify token exists, is not expired, and is not used
6. Load associated user

### Execute

1. Verify user account is still active
2. Hash new_password using bcrypt with cost factor 12
3. Update user's password hash
4. Mark reset token as used
5. Invalidate all existing sessions for this user
6. Create audit log entry
7. Send confirmation email about password change
8. Optionally: auto-login user and return tokens

### Side Effects

- [ ] Emit event: `PasswordReset`
- [ ] Mark token as used
- [ ] Invalidate all sessions
- [ ] Send confirmation email
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Missing fields or password mismatch |
| ErrInvalidToken | 400 | Token not found or malformed |
| ErrTokenExpired | 400 | Reset token has expired |
| ErrTokenUsed | 400 | Reset token already used |
| ErrWeakPassword | 400 | Password doesn't meet requirements |
| ErrAccountInactive | 403 | User account is deactivated |

## Business Rules

1. Token can only be used once
2. Token expires after 1 hour from request
3. All existing sessions must be invalidated
4. New password must meet complexity requirements
5. User should be notified via email of the change
6. If account was locked due to failed attempts, unlock it

## Notes

- Use constant-time comparison for token verification
- Consider rate limiting token verification attempts
- Log all reset attempts (success and failure)
- Include IP and device info in notification email
- Consider requiring user to log in again vs auto-login
- Delete or archive used tokens after successful reset
