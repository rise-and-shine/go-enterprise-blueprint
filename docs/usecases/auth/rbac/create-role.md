# Create Role

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.create-role

## Summary

Creates a new role in the system. Roles are containers for permissions that can be assigned to actors (users, groups, or service accounts).

## Actor

- **Primary:** Admin
- **Authorization:** `auth:role:create` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| name | string | yes | min:2, max:50, alphanumeric with hyphens, unique |
| description | string | no | max:500 |
| permissions | []string | no | array of valid permission strings |

## Output

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | Unique identifier of created role |
| name | string | Role name |
| description | string | Role description |
| permissions | []string | List of permissions assigned to role |
| created_at | timestamp | Creation timestamp |

## Flow

### Validate

1. Validate input format (name format, description length)
2. Check authorization - caller must have `auth:role:create` permission
3. Verify role name is unique
4. Validate all permissions in the list are valid system permissions

### Execute

1. Create Role record
2. If permissions provided, create RolePermission records for each
3. Create audit log entry

### Side Effects

- [ ] Emit event: `RoleCreated`
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid name format or description too long |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing `auth:role:create` permission |
| ErrConflict | 409 | Role name already exists |
| ErrInvalidPermission | 400 | One or more permissions are invalid |

## Business Rules

1. Role names must be unique across the system
2. Role names should follow kebab-case convention (e.g., "user-admin")
3. Permissions must be valid system-defined permissions
4. Cannot create roles with reserved names (e.g., "superuser", "system")
5. Empty permissions list is allowed (permissions can be added later)

## Notes

- Consider implementing role templates for common patterns
- Role name is immutable after creation
- Implement permission validation against a central permission registry
- Consider hierarchical roles in future iterations
