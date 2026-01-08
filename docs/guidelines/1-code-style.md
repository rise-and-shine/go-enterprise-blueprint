# Code Style

## Package Naming

- Lowercase, single-word, matches directory name
- **Exception:** underscore allowed for use case layer (`create_user`)

```go
// Good
package user
package auth
package create_user  // use case layer only

// Bad
package userUtils
package common_helpers
```

## File Naming

Use snake_case: `user_repository.go`, `create_order.go`

## Struct Design

- Exported fields first, then unexported
- Group related fields
- Don't repeat default tags (e.g., if `bun` defaults to lowercase, omit tag)

## Comments

- **Don't** write obvious comments (`// GetUser gets a user`)
- Start exported item comments with the item name
- Keep comments up-to-date with the code

## Formatting and Linting

- Use `make fmt` and `make lint` to format and lint
- Don't push unformatted code or code with linting errors
