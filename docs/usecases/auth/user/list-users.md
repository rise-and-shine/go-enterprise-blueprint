# List Users

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.list-users

## Summary

Retrieves a paginated list of users with optional filtering and sorting. Used for user management interfaces and reporting.

## Actor

- **Primary:** Admin
- **Authorization:** `auth:user:read` permission required

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| page | int | no | min:1, default:1 |
| page_size | int | no | min:1, max:100, default:20 |
| sort_by | string | no | enum: id, username, created_at, updated_at |
| sort_order | string | no | enum: asc, desc, default: desc |
| filter.is_active | bool | no | filter by active status |
| filter.username | string | no | partial match (contains) |
| filter.created_after | timestamp | no | ISO8601 format |
| filter.created_before | timestamp | no | ISO8601 format |

## Output

| Field | Type | Description |
|-------|------|-------------|
| users | []User | Array of user objects |
| users[].id | int64 | User ID |
| users[].username | string | Username |
| users[].person_id | int64 | Reference to Person |
| users[].is_active | bool | Active status |
| users[].created_at | timestamp | Creation timestamp |
| users[].updated_at | timestamp | Last update timestamp |
| pagination.page | int | Current page number |
| pagination.page_size | int | Items per page |
| pagination.total_items | int | Total number of matching users |
| pagination.total_pages | int | Total number of pages |

## Flow

### Validate

1. Validate pagination parameters are within bounds
2. Check authorization - caller must have `auth:user:read` permission
3. Validate sort_by is an allowed field
4. Validate date filters are valid timestamps

### Execute

1. Build query with filters
2. Apply sorting
3. Execute count query for total items
4. Execute paginated select query
5. Map results to output DTOs
6. Calculate pagination metadata

### Side Effects

- [ ] Create audit log (optional, for sensitive environments)

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid pagination or filter parameters |
| ErrUnauthorized | 401 | Not authenticated |
| ErrForbidden | 403 | Missing `auth:user:read` permission |

## Business Rules

1. Results are paginated to prevent large data transfers
2. Maximum page size is enforced (100)
3. Password hashes are never included in response
4. Sensitive fields may be masked based on caller's permission level
5. Default sort is by created_at descending (newest first)

## Notes

- Consider caching for frequently accessed pages
- Implement cursor-based pagination for large datasets
- Add export functionality for compliance reporting
- Consider field-level permissions for sensitive user data
