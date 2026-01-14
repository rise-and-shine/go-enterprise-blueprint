# Get Actor Permissions

Retrieves all effective permissions for an actor (from roles + direct permissions).

> **type**: user_action

> **operation-id**: `get-actor-permissions`

> **access**: GET /auth/get-actor-permissions

> **actor**: admin

> **permissions**: `superadmin`

## Input

Query parameters:

- `actor_type`: string, required, one of: user, admin, service_acc
- `actor_id`: string, required, UUID format

## Output

```json
{
    "actor_type": "admin",
    "actor_id": "uuid-string",
    "permissions": {
        "from_roles": ["users:read", "users:write"],
        "direct": ["system:admin"],
        "effective": ["users:read", "users:write", "system:admin"]
    }
}
```

## Execute

- Validate actor_type is valid enum
- Validate actor_id format
- Query actor's direct permissions
- Query actor's roles and their permissions
- Merge and deduplicate permissions
- Return categorized permissions

## Error Scenarios

- `INVALID_ACTOR_TYPE`: Actor type is not valid
- `INVALID_ACTOR_ID`: Actor ID format is invalid
