# Get Actor Roles

Retrieves all roles assigned to a specific actor.

> **type**: user_action

> **operation-id**: `get-actor-roles`

> **access**: GET /auth/get-actor-roles

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
    "roles": [
        {"id": 123, "name": "editor"},
        {"id": 456, "name": "viewer"}
    ]
}
```

## Execute

- Validate actor_type is valid enum
- Validate actor_id format
- Query actor roles with role details
- Return actor roles list

## Error Scenarios

- `INVALID_ACTOR_TYPE`: Actor type is not valid
- `INVALID_ACTOR_ID`: Actor ID format is invalid
