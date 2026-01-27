# Get Roles

Retrieves all roles in the system.

> **type**: user_action

> **operation-id**: `get-roles`

> **access**: GET /auth/get-roles

> **actor**: admin

> **permissions**: `superadmin`

## Input

Query parameters:

- `page`: int, optional, default 1, min 1
- `page_size`: int, optional, default 20, min 1, max 100

## Output

```json
{
    "items": [
        {
            "id": 123,
            "name": "string",
            "created_at": "2024-01-01T00:00:00Z",
            "updated_at": "2024-01-01T00:00:00Z"
        }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
}
```

## Execute

- Validate pagination parameters
- Query roles with pagination
- Return paginated list of roles
