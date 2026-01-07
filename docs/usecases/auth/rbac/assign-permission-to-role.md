# Assign Permission to Role

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.assign-permission-to-role

## Summary

Adds a permission to a role. All actors with this role will gain the new permission. Can also be used to remove permissions from a role.

## Actor

- **Primary:** Admin
- **Authorization:** `auth:permission:assign` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| role_id | int64 | yes | must exist |
| permission | string | yes | valid permission string |
| action | string | no | enum: add, remove, default: add |

## Output

| Field | Type | Description |
|-------|------|-------------|
| role_id | int64 | Role ID |
| role_name | string | Role name |
| permission | string | The permission that was modified |
| action | string | Action performed (add/remove) |
| actors_affected | int | Number of actors affected by this change |
| current_permissions | []string | Updated list of role permissions |

## Flow

### Validate

1. Validate input format
2. Check authorization - caller must have `auth:permission:assign` permission
3. Verify role exists and is not system-protected
4. Validate permission string format and existence
5. For add: verify permission not already assigned
6. For remove: verify permission is currently assigned

### Execute

1. Load Role record
2. Verify role is not system-protected
3. If action is "add":
   - Create RolePermission record
4. If action is "remove":
   - Delete RolePermission record
5. Count actors currently assigned this role
6. Invalidate authorization cache for all affected actors
7. Create audit log entry

### Side Effects

- [ ] Emit event: `RolePermissionChanged`
- [ ] Invalidate authorization cache for affected actors
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid permission format or action |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing permission or protected role |
| ErrNotFound | 404 | Role not found |
| ErrConflict | 409 | Permission already assigned (add) or not assigned (remove) |
| ErrInvalidPermission | 400 | Permission doesn't exist in system |

## Business Rules

1. Permissions must be valid system-defined permissions
2. System-protected roles cannot be modified
3. Changes take effect immediately for all role holders
4. Permission format: `module:resource:action` (e.g., `auth:user:create`)
5. Cannot remove all permissions from certain essential roles

## Notes

- Implement permission registry for validation
- Batch cache invalidation for efficiency
- Consider permission dependencies (some permissions imply others)
- Log the scope of impact (number of affected actors)
- Consider implementing permission groups for easier management
