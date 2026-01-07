# List Roles

**UseCase type:** user_action

**Module:** auth

**Operation ID:** auth.list-roles

## Summary

Retrieves a paginated list of roles with optional filtering. Used for role management interfaces and permission assignment workflows.

## Actor

- **Primary:** Admin, User (limited view)
- **Authorization:** `auth:role:read` permission for full list, limited view for authenticated users

## Input

| Field | Type | Required | Validation |
|-------|------|----------|------------|
| page | int | no | min:1, default:1 |
| page_size | int | no | min:1, max:100, default:20 |
| sort_by | string | no | enum: id, name, created_at |
| sort_order | string | no | enum: asc, desc, default: asc |
| filter.name | string | no | partial match (contains) |
| filter.include_system | bool | no | include system-protected roles, default: false |

## Output

| Field | Type | Description |
|-------|------|-------------|
| roles | []Role | Array of role objects |
| roles[].id | int64 | Role ID |
| roles[].name | string | Role name |
| roles[].description | string | Role description |
| roles[].is_system | bool | Whether this is a system-protected role |
| roles[].permissions_count | int | Number of permissions in role |
| roles[].actors_count | int | Number of actors with this role |
| roles[].created_at | timestamp | Creation timestamp |
| pagination.page | int | Current page number |
| pagination.page_size | int | Items per page |
| pagination.total_items | int | Total number of matching roles |
| pagination.total_pages | int | Total number of pages |

## Flow

### Validate

1. Validate pagination parameters
2. Check authorization level:
   - `auth:role:read` for full access
   - Authenticated users see limited public role list
3. Validate sort_by is an allowed field

### Execute

1. Build query with filters
2. Apply sorting (default: name ascending)
3. Execute count query for total items
4. Execute paginated select query
5. For each role, fetch permission count and actor count
6. Map results to output DTOs
7. Calculate pagination metadata

### Side Effects

None (read-only operation)

## Error Scenarios

| Error | Code | When |
|-------|------|------|
| ErrInvalidInput | 400 | Invalid pagination or filter parameters |
| ErrUnauthorized | 401 | Not authenticated |

## Business Rules

1. Results are paginated to prevent large data transfers
2. Default sort is by name ascending (alphabetical)
3. System roles are hidden by default
4. Non-admin users see a limited subset of roles
5. Permission and actor counts help with role management

## Notes

- Consider caching role list as it changes infrequently
- Include metadata about role usage for management decisions
- Implement search/autocomplete for role selection UI
