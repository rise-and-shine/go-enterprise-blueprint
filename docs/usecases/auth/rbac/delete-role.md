# Delete Role

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.delete-role

## Summary

Deletes a role from the system. This removes the role and all its permission assignments. Actors with this role will lose associated permissions.

## Actor

- **Primary:** Admin
- **Authorization:** `auth:role:delete` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| role_id | int64 | yes | must exist |
| force | bool | no | if true, delete even if role is assigned to actors |

## Output

| Field | Type | Description |
|-------|------|-------------|
| success | bool | True on successful deletion |
| name | string | Name of deleted role |
| actors_affected | int | Number of actors who lost this role |

## Flow

### Validate

1. Validate role_id exists
2. Check authorization - caller must have `auth:role:delete` permission
3. Verify role is not a system-protected role
4. If force is false, verify role is not assigned to any actors

### Execute

1. Load Role record
2. Count actors currently assigned this role
3. If actors exist and force is false, return error
4. Delete all ActorRole records referencing this role
5. Delete all RolePermission records for this role
6. Delete Role record
7. Create audit log entry with list of affected actors

### Side Effects

- [ ] Emit event: `RoleDeleted`
- [ ] Remove ActorRole assignments
- [ ] Remove RolePermission records
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid role_id |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing permission or protected role |
| ErrNotFound | 404 | Role not found |
| ErrRoleInUse | 400 | Role assigned to actors and force=false |

## Business Rules

1. System-protected roles cannot be deleted
2. By default, cannot delete roles that are assigned to actors
3. Force flag allows deletion even when role is in use
4. Deletion is permanent - consider soft delete for audit
5. All permission grants through this role are immediately revoked

## Notes

- Consider implementing soft delete for compliance requirements
- Log all actors affected by the deletion
- Consider notification to affected users
- Cache invalidation may be required for authorization checks
- This is a dangerous operation - require confirmation in UI
