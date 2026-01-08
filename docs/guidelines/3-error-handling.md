# Error Handling

Use `github.com/code19m/errx` package for all error handling.

## Wrapping Errors

Always wrap errors when returning:

```go
if err := doSomething(); err != nil {
    return errx.Wrap(err)
}
```

## Error Codes

Define codes instead of sentinel errors:

```go
var (
    CodeNotFound     = "NOT_FOUND"
    CodeUnauthorized = "UNAUTHORIZED"
    CodeInvalidInput = "INVALID_INPUT"
)
```

## Checking Errors

Use `errx.IsCodeIn` for error-based logic:

```go
if errx.IsCodeIn(err, CodeNotFound) {
    // handle not found
}
```

## Error Types

- Error types are defined **only at use case layer**
- All downstream errors should return `errx.T_Internal` (default)
- Use case layer knows its caller (actor) and assigns appropriate types
- Use `errx.WrapWithTypeOnCodes` to change type based on specific codes (e.g., when error is related to user input)
