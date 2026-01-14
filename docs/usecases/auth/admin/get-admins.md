# Get Admins

Retrieves a paginated list of all admin accounts.

> **type**: user_action

> **operation-id**: `get-admins`

> **access**: GET /auth/get-admins

> **actor**: admin

> **permissions**: `superadmin`

## Input

Query parameters:

- `page`: int, optional, default 1, min 1
- `page_size`: int, optional, default 20, min 1, max 100
- `is_active`: bool, optional, filter by active status

## Output

```json
{
    "items": [
        {
            "id": "uuid-string",
            "username": "string",
            "is_superadmin": false,
            "is_active": true,
            "last_active_at": "2024-01-01T00:00:00Z",
            "created_at": "2024-01-01T00:00:00Z"
        }
    ],
    "total": 50,
    "page": 1,
    "page_size": 20
}
```

## Execute

- Validate pagination parameters
- Apply filters if provided
- Query admins with pagination
- Return paginated list of admins (without password hashes)
