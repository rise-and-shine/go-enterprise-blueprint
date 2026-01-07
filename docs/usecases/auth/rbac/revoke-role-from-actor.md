# Revoke Role from Actor

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.revoke-role-from-actor

## Summary

Removes a role assignment from an actor. The actor will lose all permissions associated with the role (unless granted by another role or directly).

## Actor

- **Primary:** Admin
- **Authorization:** `auth:role:revoke` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| role_id | int64 | yes | must exist |
| actor_type | string | yes | enum: user, group, service_acc |
| actor_id | string | yes | must have the role assigned |

## Output

| Field | Type | Description |
|-------|------|-------------|
| success | bool | True on successful revocation |
| role_name | string | Name of revoked role |
| actor_type | string | Type of actor |
| actor_id | string | ID of the actor |
| permissions_revoked | []string | Permissions no longer available (unless from other source) |

## Flow

### Validate

1. Validate input format
2. Check authorization - caller must have `auth:role:revoke` permission
3. Verify role exists
4. Verify actor_type is valid
5. Verify actor has this role assigned
6. Check if revoking would leave system without required roles

### Execute

1. Load ActorRole assignment
2. Verify this isn't the last superuser role removal
3. Delete ActorRole record
4. Calculate which permissions are actually lost (not granted elsewhere)
5. Invalidate authorization cache for this actor
6. Create audit log entry

### Side Effects

- [ ] Emit event: `RoleRevoked`
- [ ] Invalidate authorization cache
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid actor_type or missing fields |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing `auth:role:revoke` permission |
| ErrNotFound | 404 | Role or assignment not found |
| ErrLastSuperuser | 400 | Cannot revoke last superuser role |

## Business Rules

1. Cannot revoke role that isn't assigned
2. System must always have at least one superuser
3. Revocation takes effect immediately
4. Permissions may still be available from other roles or direct grants
5. Admins cannot revoke their own essential roles (self-lockout prevention)

## Notes

- Calculate effective permission change for user feedback
- Cache invalidation is critical for security
- Consider soft-revoke with grace period for sensitive roles
- Log detailed information for audit compliance
- Consider notification to affected user
