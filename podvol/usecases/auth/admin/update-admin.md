# Update Admin

Updates an existing admin's information, optionally including password change.

> **type**: user_action

> **operation-id**: `update-admin`

> **access**: POST /auth/update-admin

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
  "id": "uuid-string", // required
  "username": "string", // optional, 3-50 chars, unique if provided
  "password": "string", // optional, min 8 chars if provided
  "is_superadmin": false // optional
}
```

## Output

```json
{
  "id": "uuid-string",
  "username": "string",
  "is_superadmin": false,
  "is_active": true,
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Execute

- Validate admin ID exists
- If username provided, check uniqueness (excluding current admin)
- If password provided, validate strength and hash it
- Update admin record with provided fields
- Return updated admin (without password hash)

## Error Scenarios

- `ADMIN_NOT_FOUND`: Admin does not exist
- `USERNAME_EXISTS`: Another admin with this username already exists
- `INVALID_USERNAME`: Username does not meet requirements
- `WEAK_PASSWORD`: Password does not meet minimum strength requirements
- `CANNOT_DEMOTE_LAST_SUPERADMIN`: Cannot remove superadmin status from the last superadmin
