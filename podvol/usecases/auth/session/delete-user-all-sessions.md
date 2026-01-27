# Delete User All Sessions

Deletes all sessions for a specific actor. Only superadmins can perform this operation.

> **type**: user_action

> **operation-id**: `delete-user-all-sessions`

> **access**: POST /auth/delete-user-all-sessions

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "actor_type": "string", // required, one of: user, admin, service_acc
    "actor_id": "string"    // required, UUID format
}
```

## Output

```json
{
    "deleted_count": 5
}
```

## Execute

- Validate actor_type is valid enum value
- Validate actor_id is valid UUID
- Delete all sessions matching actor_type and actor_id
- Return count of deleted sessions

## Error Scenarios

- `INVALID_ACTOR_TYPE`: Actor type is not valid
- `INVALID_ACTOR_ID`: Actor ID is not valid UUID format
