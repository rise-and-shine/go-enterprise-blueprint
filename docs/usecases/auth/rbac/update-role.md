# Update Role

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.update-role

## Summary

Updates an existing role's metadata. This does not modify permissions - use assign-permission-to-role for that.

## Actor

- **Primary:** Admin
- **Authorization:** `auth:role:update` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| role_id | int64 | yes | must exist |
| description | string | no | max:500 |

## Output

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | Role ID |
| name | string | Role name (unchanged) |
| description | string | Updated description |
| updated_at | timestamp | Update timestamp |

## Flow

### Validate

1. Validate role_id exists
2. Check authorization - caller must have `auth:role:update` permission
3. Verify role is not a system-protected role
4. Validate description length

### Execute

1. Load Role record
2. Update description field
3. Update timestamp
4. Save Role record
5. Create audit log entry

### Side Effects

- [ ] Emit event: `RoleUpdated`
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Description too long |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing permission or protected role |
| ErrNotFound | 404 | Role not found |

## Business Rules

1. Role name cannot be changed (immutable)
2. System-protected roles (e.g., "superuser") cannot be modified
3. Only description can be updated via this use case
4. Permission changes require separate use case

## Notes

- Role name is intentionally immutable to maintain referential integrity
- System roles should be marked with a flag to prevent modification
- Consider versioning for role changes in compliance-heavy environments
