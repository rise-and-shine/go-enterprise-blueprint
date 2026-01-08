# Testing

## Mocking

Use `mockery` v2 to generate mocks. Don't write mocks by hand.

```bash
mockery --name=UserRepository --output=mocks
```

## Fake Objects

For external dependencies (APIs, services), write test doubles/fakes for integration tests:

```go
type FakeUserService struct {
    users map[int64]*User
}

func (f *FakeUserService) GetUser(ctx context.Context, id int64) (*User, error) {
    return f.users[id], nil
}
```

## Race Detection

Always run tests with race detector:

```bash
go test -race ./...
```
