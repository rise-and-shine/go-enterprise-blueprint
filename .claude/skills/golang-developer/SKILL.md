---
name: golang-developer
description: Ensures Go code follows idiomatic patterns, best practices, community and team standards. Use when writing Go code, reviewing Go implementations, or refactoring Go projects.
---

# Golang Developer Skill

## When to Use

- Writing new Go code
- Reviewing Go code
- Refactoring Go code

## Guidelines

Scan [docs/guidelines/](../../../docs/guidelines/) for full team guidelines:

---

## Critical Reminders

These are the most common mistakes. Always check before submitting code.

### Context Parameter

Every function that performs I/O or calls other services MUST accept `ctx context.Context` as the first parameter. Never create context inside a function.

```go
// Wrong
func FetchUser(id int64) (*User, error) {
    ctx := context.Background()
    return db.QueryUser(ctx, id)
}

// Correct
func FetchUser(ctx context.Context, id int64) (*User, error) {
    return db.QueryUser(ctx, id)
}
```

### Error Handling

Never ignore errors. Every error must be explicitly handled or wrapped with `errx.Wrap(err)`.

```go
// Wrong
doSomething()
_ = doSomething()

// Correct
if err := doSomething(); err != nil {
    return errx.Wrap(err)
}
```

### Comments

Never write comments describing your actions or thought process. Comments like `// I fixed this`, `// Added this to handle...`, or `// This was missing` are prohibited.
