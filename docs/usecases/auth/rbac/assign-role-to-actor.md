# Assign Role to Actor

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.assign-role-to-actor

## Summary

Assigns a role to an actor (user, group, or service account). The actor will gain all permissions associated with the role.

## Actor

- **Primary:** Admin
- **Authorization:** `auth:role:assign` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| role_id | int64 | yes | must exist |
| actor_type | string | yes | enum: user, admin, service_acc |
| actor_id | string | yes | must reference existing actor of specified type |

## Output

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | ActorRole assignment ID |
| role_id | int64 | Assigned role ID |
| role_name | string | Assigned role name |
| actor_type | string | Type of actor |
| actor_id | string | ID of the actor |
| permissions_granted | []string | List of permissions now available to actor |
| created_at | timestamp | Assignment timestamp |

## Flow

### Validate

1. Validate input format
2. Check authorization - caller must have `auth:role:assign` permission
3. Verify role exists
4. Verify actor_type is valid
5. Verify actor exists (user/admin/service account)
6. Check if assignment already exists

### Execute

1. Load Role record
2. Verify actor exists based on actor_type
3. Check for duplicate assignment
4. Create ActorRole record
5. Fetch role's permissions for response
6. Invalidate authorization cache for this actor
7. Create audit log entry

### Side Effects

- [ ] Emit event: `RoleAssigned`
- [ ] Invalidate authorization cache
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid actor_type or missing fields |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing `auth:role:assign` permission |
| ErrNotFound | 404 | Role or actor not found |
| ErrConflict | 409 | Role already assigned to this actor |

## Business Rules

1. Same role cannot be assigned twice to the same actor
2. Actor must exist before role can be assigned
3. Assignment takes effect immediately
4. Group role assignments cascade to all group members
5. Service account roles are used for API authorization
6. Cannot assign system-exclusive roles to regular users

## Notes

- Consider implementing role expiration (time-limited assignments)
- Cache invalidation is critical for authorization consistency
- Log who assigned the role for audit trail
- Consider notification to user when roles change
- Implement maximum roles per actor limit if needed
