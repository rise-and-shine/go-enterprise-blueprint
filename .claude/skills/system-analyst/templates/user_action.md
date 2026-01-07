# {Use Case Name}

{One-two sentence summary of what this does and why.}

> **type**: user_action

> **operation-id**: `{operation-id}`

> **access**: {GET or POST} {path}

> **actor**: {one of user, admin, service_acc}

> **permissions**: `{permissions list}`

## Input

EXAMPLE FOR POST

```json
{
    "username": "string" // required, min=5, max=50
    "email": "string" // required, email format
    "password": "string" // required, min=8
    "age": 12 // optional, integer
    "date_of_birth": "2001-01-31" // optional, date format
}
```

EXAMPLE FOR GET

- `username`: string, required, 3-50 chars
- `email`: string, required, email format
- `is_active`: bool, optional, default true

## Output

```json
{
    "id": "string",
    "username": "string",
    "email": "string",
    "is_active": bool,
    "age": 12, // nullable
    "date_of_birth": "2001-01-31" // nullable
}
```

## Execute

- Validate {what}
- Check {precondition}
- Process {what}
- Produce `{EventName}` event
- Return {result}

## Special Error Codes (THIS SECTIN IS REQUIRED ONLY FOR service_acc ACTOR TYPES)

- `{error code}`: {optinal description of error}

Examples

- `EMAIL_ALREADY_EXISTS`: User with provided email already exists
