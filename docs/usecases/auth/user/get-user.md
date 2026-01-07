# Get User

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.get-user

## Summary

Retrieves detailed information about a specific user by their ID. Used for viewing user profiles and management details.

## Actor

- **Primary:** Admin, User (self)
- **Authorization:** `auth:user:read` permission OR viewing own profile

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| user_id | int64 | yes | must exist |

## Output

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | User ID |
| username | string | Username |
| person_id | int64 | Reference to Person |
| person | Person | Embedded Person details (if authorized) |
| person.pin | string | Person's PIN (masked for non-admins) |
| person.birth_date | string | Birth date |
| person.birth_place | string | Birth place |
| is_active | bool | Active status |
| roles | []Role | Assigned roles |
| permissions | []string | Direct permissions |
| created_at | timestamp | Creation timestamp |
| updated_at | timestamp | Last update timestamp |
| last_login_at | timestamp | Last successful login (if tracked) |

## Flow

### Validate

1. Validate user_id format
2. Check authorization:
   - Admin with `auth:user:read` can view any user
   - Users can view their own profile
3. Verify user exists

### Execute

1. Load User record by ID
2. Load associated Person record
3. Load user's roles via ActorRole
4. Load user's direct permissions via ActorPermission
5. Apply field masking based on caller's permissions
6. Return assembled user details

### Side Effects

- [ ] Create audit log (optional, for sensitive data access)

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid user_id format |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing permission and not own profile |
| ErrNotFound | 404 | User not found |

## Business Rules

1. Password hash is never returned
2. Users can always view their own profile
3. Person's PIN should be masked for non-admin viewers
4. Sensitive fields may have different visibility levels
5. Include role and permission information for authorization context

## Notes

- Consider implementing field-level access control
- Cache frequently accessed user profiles
- Person's full details might require additional permissions
- Consider rate limiting to prevent enumeration attacks
