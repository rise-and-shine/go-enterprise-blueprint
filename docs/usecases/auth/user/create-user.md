# Create User

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.create-user

## Summary

Creates a new user account in the system. This is the standard user registration flow that creates a User entity linked to a Person record.

## Actor

- **Primary:** Admin, System
- **Authorization:** `auth:user:create` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| username | string | yes | min:3, max:50, alphanumeric with underscore, unique |
| password | string | yes | min:8, max:128, must contain uppercase, lowercase, digit |
| person_id | int64 | yes | must reference existing Person |
| is_active | bool | no | defaults to true |

## Output

| Field | Type | Description |
|-------|------|-------------|
| id | int64 | Unique identifier of created user |
| username | string | Username of the created user |
| person_id | int64 | Reference to associated Person |
| is_active | bool | Account active status |
| created_at | timestamp | Creation timestamp |

## Flow

### Validate

1. Validate input format and constraints (username format, password strength)
2. Check authorization - caller must have `auth:user:create` permission
3. Verify username is unique in the system
4. Verify person_id references an existing Person record
5. Verify Person is not already linked to another User

### Execute

1. Hash the password using bcrypt with cost factor 12
2. Create User record with provided data
3. Link User to Person record
4. Create audit log entry for user creation

### Side Effects

- [ ] Emit event: `UserCreated`
- [ ] Create audit log

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Username format invalid, password too weak |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing `auth:user:create` permission |
| ErrNotFound | 404 | Person with given ID not found |
| ErrConflict | 409 | Username already exists or Person already linked to User |

## Business Rules

1. Username must be unique across the entire system
2. One Person can only be linked to one User account
3. Password must meet minimum security requirements
4. New users are active by default unless explicitly set otherwise

## Notes

- Password is never stored in plain text, only bcrypt hash
- Consider rate limiting this endpoint to prevent abuse
- Username cannot be changed after creation (use a separate use case if needed)
