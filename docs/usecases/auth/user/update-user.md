# Update User

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.update-user

## Summary

Updates an existing user's profile information. This allows modification of user attributes except for immutable fields like username and ID.

## Actor

- **Primary:** Admin, User (self)
- **Authorization:** `auth:user:update` permission OR updating own profile

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| user_id | int64 | yes | must exist |
| is_active | bool | no | - |
| person_id | int64 | no | must reference existing Person, not linked to another User |

## Output

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | User ID |
| username | string | Username (unchanged) |
| person_id | int64 | Reference to associated Person |
| is_active | bool | Updated active status |
| updated_at | timestamp | Update timestamp |

## Flow

### Validate

1. Validate user_id exists in the system
2. Check authorization:
   - Admin with `auth:user:update` can update any user
   - Users can update their own profile (limited fields)
3. If person_id provided, verify Person exists and is not linked to another User
4. If is_active is being set to false, verify caller is not deactivating themselves

### Execute

1. Load current User record
2. Apply updates to allowed fields
3. Update the updated_at timestamp
4. Save User record
5. Create audit log entry

### Side Effects

- [ ] Emit event: `UserUpdated`
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid field values |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing permission and not own profile |
| ErrNotFound | 404 | User or Person not found |
| ErrConflict | 409 | Person already linked to another User |
| ErrSelfDeactivation | 400 | Attempting to deactivate own account |

## Business Rules

1. Username cannot be changed (immutable after creation)
2. User ID cannot be changed
3. Users cannot deactivate their own account
4. Only admins can change is_active status
5. Person linkage can only be changed by admins

## Notes

- Consider implementing field-level permissions for fine-grained access control
- Changes to is_active may trigger session invalidation
- Audit log should capture before/after values for compliance
