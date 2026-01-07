# Change Password

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.change-password

## Summary

Allows an authenticated user to change their own password. Requires the current password for verification before setting the new password.

## Actor

- **Primary:** Authenticated User
- **Authorization:** Valid session, changing own password only

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| current_password | string | yes | min:1, must match current |
| new_password | string | yes | min:8, max:128, complexity requirements |
| confirm_password | string | yes | must match new_password |

## Output

| Field | Type | Description |
|-------|------|-------------|
| success | bool | Always true on success |
| message | string | Confirmation message |
| sessions_invalidated | int | Number of other sessions terminated |

## Flow

### Validate

1. Validate all password fields are present
2. Verify new_password matches confirm_password
3. Validate new password meets complexity requirements:
   - Minimum 8 characters
   - At least one uppercase letter
   - At least one lowercase letter
   - At least one digit
   - At least one special character (recommended)
4. Verify new password is different from current password
5. Check password is not in common password list
6. Verify current_password matches stored hash

### Execute

1. Load current user's record
2. Verify current_password against stored hash
3. Hash new_password using bcrypt with cost factor 12
4. Update user's password hash
5. Invalidate all other sessions (security best practice)
6. Keep current session active
7. Create audit log entry
8. Send email notification about password change

### Side Effects

- [ ] Emit event: `PasswordChanged`
- [ ] Invalidate other sessions
- [ ] Send notification email
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Missing fields or password mismatch |
| ErrUnauthorized | 401 | Not authenticated |
| ErrWrongPassword | 401 | Current password is incorrect |
| ErrWeakPassword | 400 | New password doesn't meet requirements |
| ErrSamePassword | 400 | New password same as current |
| ErrCommonPassword | 400 | Password is in common password list |

## Business Rules

1. User must verify identity with current password
2. New password must meet complexity requirements
3. New password cannot be the same as current password
4. Password history may be checked (last N passwords)
5. All other sessions should be invalidated for security
6. User should be notified via email of the change
7. Rate limit password change attempts

## Notes

- Consider implementing password history (prevent reuse of last 5 passwords)
- Send notification to all registered email addresses
- Include IP and device info in notification email
- Consider requiring re-authentication for sensitive accounts
- Log password changes for security audit
