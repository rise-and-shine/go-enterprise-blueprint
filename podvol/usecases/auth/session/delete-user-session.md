# Delete User Session

Deletes a specific session for an actor. Superadmins can delete any session, while users with `session:delete:own` permission can only delete their own sessions.

> **type**: user_action

> **operation-id**: `delete-user-session`

> **access**: POST /auth/delete-user-session

> **actor**: admin

> **permissions**: `superadmin` OR `session:delete:own` (for own sessions only)

## Input

```json
{
    "session_id": 123 // required, int64
}
```

## Output

```json
{
    "success": true
}
```

## Execute

- Validate session_id
- Find session by ID
- If actor is not superadmin, verify session belongs to the actor (same actor_type and actor_id)
- Delete the session
- Return success

## Error Scenarios

- `SESSION_NOT_FOUND`: Session does not exist
- `PERMISSION_DENIED`: Actor is not authorized to delete this session
