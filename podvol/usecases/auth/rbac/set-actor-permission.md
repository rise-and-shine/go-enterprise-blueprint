# Set Actor Permission

Assigns direct permissions to an actor (bypassing roles). Replaces all existing direct permissions.

> **type**: user_action

> **operation-id**: `set-actor-permission`

> **access**: POST /auth/set-actor-permission

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "actor_type": "string",      // required, one of: user, admin, service_acc
    "actor_id": "string",        // required, UUID format
    "permissions": ["string"]    // required, array of permission strings
}
```

## Output

```json
{
    "actor_type": "admin",
    "actor_id": "uuid-string",
    "permissions": ["users:delete", "system:admin"]
}
```

## Execute

- Validate actor_type is valid enum
- Validate actor_id format
- Validate permission strings format
- Delete existing actor permissions
- Insert new actor permissions
- Return updated actor permissions

## Error Scenarios

- `INVALID_ACTOR_TYPE`: Actor type is not valid
- `INVALID_PERMISSION_FORMAT`: Permission string format is invalid
