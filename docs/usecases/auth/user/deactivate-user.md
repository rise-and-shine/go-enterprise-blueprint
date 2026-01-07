# Deactivate User

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.deactivate-user

## Summary

Deactivates a user account, preventing them from authenticating while preserving their data for audit and potential reactivation. This is a soft-delete operation.

## Actor

- **Primary:** Admin
- **Authorization:** `auth:user:deactivate` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| user_id | int64 | yes | must exist, must be currently active |
| reason | string | no | max:500, reason for deactivation |

## Output

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | User ID |
| username | string | Username of deactivated user |
| is_active | bool | Will be false |
| deactivated_at | timestamp | Deactivation timestamp |
| deactivated_by | int64 | ID of admin who performed deactivation |

## Flow

### Validate

1. Validate user_id exists in the system
2. Check authorization - caller must have `auth:user:deactivate` permission
3. Verify target user is currently active
4. Verify caller is not deactivating themselves
5. Verify target user is not the last active superuser

### Execute

1. Load User record
2. Set is_active to false
3. Record deactivation metadata (timestamp, actor, reason)
4. Invalidate all active sessions for this user
5. Revoke all active refresh tokens
6. Save User record
7. Create audit log entry

### Side Effects

- [ ] Emit event: `UserDeactivated`
- [ ] Invalidate user sessions
- [ ] Revoke refresh tokens
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Reason exceeds max length |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing `auth:user:deactivate` permission |
| ErrNotFound | 404 | User not found |
| ErrAlreadyDeactivated | 400 | User is already inactive |
| ErrSelfDeactivation | 400 | Attempting to deactivate own account |
| ErrLastSuperuser | 400 | Cannot deactivate the last active superuser |

## Business Rules

1. Deactivated users cannot log in
2. User data is preserved (soft delete)
3. All active sessions must be terminated immediately
4. Admins cannot deactivate themselves
5. System must always have at least one active superuser
6. Deactivation reason should be recorded for compliance

## Notes

- This is reversible - users can be reactivated via update-user
- Consider implementing a grace period before permanent data removal
- Deactivated users' data should still be accessible for audit purposes
- API tokens issued to this user should also be revoked
