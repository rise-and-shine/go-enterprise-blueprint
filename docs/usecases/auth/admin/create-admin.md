# Create Admin

Creates a new admin user account.

> **type**: user_action

> **operation-id**: `create-admin`

> **access**: POST /auth/create-admin

> **actor**: admin

> **permissions**: `superadmin`

## Input

```json
{
    "username": "string",    // required, 3-50 chars, unique
    "password": "string",    // required, min 8 chars
    "is_superadmin": false   // optional, default false
}
```

## Output

```json
{
    "id": "uuid-string",
    "username": "string",
    "is_superadmin": false,
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
}
```

## Execute

- Validate input (username, password strength)
- Check if username is unique
- Hash the password using bcrypt
- Create admin record with `is_active=true`
- Return created admin (without password hash)

## Error Scenarios

- `USERNAME_EXISTS`: Admin with this username already exists
- `INVALID_USERNAME`: Username does not meet requirements
- `WEAK_PASSWORD`: Password does not meet minimum strength requirements
