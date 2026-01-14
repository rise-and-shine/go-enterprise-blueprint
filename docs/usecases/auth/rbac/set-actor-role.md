# Set Actor Role

Assigns roles to an actor. Replaces all existing role assignments.

> **type**: user_action

> **operation-id**: `set-actor-role`

> **access**: POST /auth/set-actor-role

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "actor_type": "string",  // required, one of: user, admin, service_acc
    "actor_id": "string",    // required, UUID format
    "role_ids": [123, 456]   // required, array of role IDs
}
```

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
- Validate all role_ids exist
- Delete existing actor roles
- Insert new actor roles
- Return updated actor roles

## Error Scenarios

- `INVALID_ACTOR_TYPE`: Actor type is not valid
- `ROLE_NOT_FOUND`: One or more roles do not exist
