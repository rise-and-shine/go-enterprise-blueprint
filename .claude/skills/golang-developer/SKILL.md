---
name: golang-developer
description: Ensures Go code follows idiomatic patterns, best practices, community and team standards. Use when writing Go code, reviewing Go implementations, or refactoring Go projects.
---

# Golang Developer Skill

## When to use this Skill

Use this skill when:

- Writing new Go code
- Reviewing Go code for idioms
- Refactoring Go code
- Optimizing Go performance
- Ensuring Go code quality
- Teaching Go best practices

## Go Best Practices Checklist

### 1. Code Organization

**Package Structure**

- [ ] Package names are lowercase, single-word, concise
- [ ] Package names match directory names
- [ ] No underscores or mixedCaps in package names (underscore is allowed only for use case layer. This is teams convention)
- [ ] Avoid generic names like `util`, `common`, `base`
- [ ] Group related functionality in packages

```go
// Good package names
package user
package auth
package httpclient
package create_user // allowed only for use case layer

// Bad package names
package userUtils
package common_helpers
package pkg
```

**File Organization**

- [ ] One logical concept per file
- [ ] File names use underscores (snake_case)

### 2. Naming Conventions

**Variables**

- [ ] Use camelCase for local variables
- [ ] Use PascalCase for exported identifiers
- [ ] Short names for short scopes (`i`, `r`, `w`)
- [ ] Descriptive names for larger scopes
- [ ] Avoid stuttering (e.g., `user.UserID` → `user.ID`)

```go
// Good
var count int
var userCount int
var authenticatedUserCount int

// Bad
var user_count int
var cnt int  // unclear in large scope
```

**Functions & Methods**

- [ ] Use verbs or verb phrases
- [ ] Getters don't use `Get` prefix
- [ ] Setters use `Set` prefix
- [ ] Boolean checks use `Is`, `Has`, `Can`

```go
// Good
func (u *User) Name() string          // getter
func (u *User) SetName(name string)   // setter
func (u *User) IsActive() bool        // boolean
func (u *User) HasPermission() bool

// Bad
func (u *User) GetName() string
func (u *User) Active() bool  // unclear if getter or checker
```

**Interfaces**

- [ ] Single-method interfaces end with `-er`
- [ ] Describe capability, not data

```go
// Good
type Reader interface { Read(p []byte) (n int, err error) }
type Writer interface { Write(p []byte) (n int, err error) }
type Closer interface { Close() error }

// Bad
type IReader interface { ... }  // don't use "I" prefix
type ReaderInterface interface { ... }
```

### 3. Error Handling

**Always Handle Errors. Use github.com/code19m/errx package for error handling**

```go
// Good
if err := doSomething(); err != nil {
    return errx.Wrap(err)
}

// Bad
_ = doSomething()  // ignoring error
doSomething()      // ignoring error
```

**Custom Errors**

```go
// Though go community prefers using sentinal errors we do different approach.
// We define error codes instead of using sentinel errors
var (
    CodeNotFound      = "NOT_FOUND"
    CodeUnauthorized  = "UNAUTHORIZED"
    CodeInvalidInput  = "INVALID_INPUT"
)
```

**Error Checking**

```go
// Use errx.GetCode or errx.IsCodeIn functions for error related logic
if errx.IsCodeIn(err, CodeNotFound) {
    // handle not found
}
```

### 4. Struct Design

**Field Ordering**

- [ ] Exported fields first
- [ ] Group related fields
- [ ] Order by importance/usage
- [ ] Consider memory alignment (optional optimization)
- [ ] Don't repeat default tags (e.g., if tag `bun` defaults to fields lowercase, don't repeat it)

```go
// Good
type User struct {
    // Exported fields
    ID        int64
    Username  string
    Email     string
    CreatedAt time.Time

    // Unexported fields
    passwordHash []byte
    salt         []byte
}
```

**Struct Tags**

```go
type User struct {
    ID       int64     `json:"id"`
    Username string    `json:"username" validate:"required,min=3"`
    Email    string    `json:"email" validate:"required,email"`
    FirstName string   `json:"first_name"`
}
```

**Constructor Functions**

```go
// Provide constructor for complex initialization
func NewUser(username, email string) *User {
    return &User{
        ID:        generateID(),
        Username:  username,
        Email:     email,
        CreatedAt: time.Now(),
    }
}
```

### 5. Concurrency

**Goroutines**

- [ ] Never start a goroutine without knowing how it will stop
- [ ] Use context for cancellation
- [ ] Avoid goroutine leaks
- [ ] Use WaitGroups or channels to coordinate

```go
// Good - controlled goroutine with context
func processItems(ctx context.Context, items []Item) error {
    for _, item := range items {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            if err := process(item); err != nil {
                return err
            }
        }
    }
    return nil
}

// Bad - leaked goroutine
func badExample() {
    go func() {
        for {
            // This goroutine never stops!
            doWork()
        }
    }()
}

// Use errgroup when appropriate
func fetchAll(ctx context.Context, urls []string) error {
    g, ctx := errgroup.WithContext(ctx)
    for _, url := range urls {
        g.Go(func() error {
            return fetch(ctx, url)
        })
    }
    return g.Wait()
}
```

**Channels**

- [ ] Channel owners close channels
- [ ] Receivers never close channels
- [ ] Don't send on closed channels
- [ ] Use buffered channels appropriately

```go
// Good - owner closes
func generate(ctx context.Context) <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)  // owner closes
        for i := 0; i < 10; i++ {
            select {
            case ch <- i:
            case <-ctx.Done():
                return
            }
        }
    }()
    return ch
}
```

**Mutexes**

- [ ] Keep critical sections small
- [ ] Don't hold locks during I/O
- [ ] Use defer for unlocking
- [ ] Prefer sync.RWMutex for read-heavy workloads

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]interface{}
}

func (c *Cache) Get(key string) (interface{}, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    val, ok := c.items[key]
    return val, ok
}

func (c *Cache) Set(key string, val interface{}) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = val
}
```

**Race Detection**

```bash
# Always test with race detector
go test -race ./...
```

### 6. Context Usage

**Always Accept Context**

- [ ] First parameter should be `ctx context.Context`
- [ ] Pass context through call chain
- [ ] Don't store context in structs
- [ ] Use context for cancellation, deadlines, values

```go
// Good
func FetchUser(ctx context.Context, id int64) (*User, error) {
    // Use context for database calls, HTTP requests, etc.
    return db.QueryUser(ctx, id)
}

// Bad
func FetchUser(id int64) (*User, error) {
    ctx := context.Background()  // Don't create context inside
    return db.QueryUser(ctx, id)
}
```

**Context Values**

- [ ] Use context values sparingly
- [ ] Only for request-scoped data
- [ ] Use typed keys to avoid collisions

```go
// Define typed key
type contextKey string

const userIDKey contextKey = "userID"

// Store value
ctx = context.WithValue(ctx, userIDKey, userID)

// Retrieve value
userID, ok := ctx.Value(userIDKey).(int64)
```

### 7. Interface Usage

**Accept Interfaces, Return Structs**

```go
// Good
func ProcessData(r io.Reader) (*Result, error) {
    // Function accepts interface
    data, err := io.ReadAll(r)
    // ...
    return &Result{...}, nil  // Returns concrete type
}

// Bad
func ProcessData(r *os.File) (*Result, error) {
    // Too specific - limits testing and reuse
}
```

### 8. Performance

**Avoid Allocations** (if possible)

```go
// Good - reuse slice capacity
results := make([]Result, 0, len(inputs))
for _, input := range inputs {
    results = append(results, process(input))
}

// Bad - reallocations
var results []Result
for _, input := range inputs {
    results = append(results, process(input))  // grows dynamically
}
```

**String Building**

```go
// Good - use strings.Builder if string manipulation is very large
var b strings.Builder
b.Grow(estimatedSize)  // optional: pre-allocate
for _, str := range strs {
    b.WriteString(str)
}
result := b.String()

// Bad - string concatenation
var result string
for _, str := range strs {
    result += str  // creates new string each iteration
}
```

**Benchmarking**

```go
func BenchmarkProcess(b *testing.B) {
    data := generateTestData()
    b.ResetTimer()  // Reset timer after setup

    for i := 0; i < b.N; i++ {
        process(data)
    }
}
```

### 9. Code Style

**Formatting and linting**

- [ ] Use commands provided by Makefile for formatting and linting

**Comments**

- [ ] Don't write comments for your actions (e.g., BAD `// I fixed this...`)
- [ ] Don't write obvious comments (e.g., BAD `// getUserByID gets a user by ID`)
- [ ] Write comments for imperative code
- [ ] Write comments for interfaces with possible return error codes
- [ ] Use complete sentences, end with period
- [ ] Start with the name of the item

### 10. Logging

Todo...

## Common Anti-Patterns to Avoid

### ❌ Ignoring Errors

Always handle errors, don't ignore them.

```go
// Don't
doSomething()
```

### ❌ Using panic in Libraries

```go
// Don't - return errors instead
if err != nil {
    panic(err)
}

// Do
if err != nil {
    return errx.Wrap(err)
}
```

### ❌ Goroutine Leaks

```go
// Don't
go func() {
    for {
        doWork()  // runs forever
    }
}()

// Do
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    for {
        select {
        case <-ctx.Done():
            return
        default:
            doWork()
        }
    }
}()
```

### ❌ Naked Returns

```go
// Don't
func split(sum int) (x, y int) {
    x = sum * 4 / 9
    y = sum - x
    return  // naked return
}

// Do
func split(sum int) (int, int) {
    x := sum * 4 / 9
    y := sum - x
    return x, y
}
```

## Backend Skills

### Entities

- [ ] For entities that represent table on postgres embed pg.BaseModel
- [ ] Use pointer values for nullability
- [ ] Don't repeat default tags

### Controllers

- [ ] Keep one to one relationship between controllers and use cases
- [ ] Single use case cannot be called from multiple controllers
- [ ] Keep this layer simple and clean
- [ ] Don't write manual controllers, use general components (like forward.ToUseCase for http controllers)

### Use Cases

Each use case has it's concrete type (user_action, event_subscriber, async_task, manual_command) and concrete signature.

- [ ] Keep document first approach, don't write use case before documenting it by our system analyst rules
- [ ] Use cases should reference to it's documentation (in comments)
- [ ] Use cases should reflect documentation as much as possible
- [ ] Always enforce input validation (if not validated on controller layer) and business logic validation rules
- [ ] If same logic duplicated in multiple use cases consider packaging those logic into PBLC layer
- [ ] Separate file for each use case
- [ ] For OperationID define constants in top of use case file like `create-superuser`

### PBLC layer (Packaged Business Logic Component)

This layer is used for deduplicating and packaging independent, reusable business logic components. These components will be called ONLY from use cases. These components don't know about their callers, so validating input values should be strict in this layer. Also this layer is most suitable for implementing OOP design patterns like State, Strategy if NEEDED. In this layer we are allowed to design our components in any way we want.

### Infrastructure/Repository layer

- [ ] Prefer Generalization over Specialization
- [ ] For postgres repositories use repogen package. Try to minimize additioal methods, try to use general-purpose methods of repogen on caller layers
- [ ] For http clients use resty v2

Examples:
postgres store - TODO...
http clients - TODO...
in-memory components - TODO...

### Transaction management

Manage atomic transactions on use case layer with UnitOfWork pattern.
Repository layer shouldn't be responsible for transaction management. It works with bun.IDB.

Examples: TODO...

## Testing

### Table-Driven Tests

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -2, -3, -5},
        {"mixed", -2, 3, 1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Add(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Add(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

### Mocking

- [ ] Use interfaces for mocking
- [ ] Use `mockery` v2 for mocking, don't write mocks by hand

### Fake Objects

- [ ] For external dependencies write test doubles/fakes for using in integration tests

```go
type UserService interface {
    GetUser(ctx context.Context, id int64) (*User, error)
}

type FakeUserService struct {
    fakeUsers map[int64]*User
}

func (f *FakeUserService) GetUser(ctx context.Context, id int64) (*User, error) {
    return f.fakeUsers[id], nil
}

func NewFakeUserService(fakeUsers map[int64]*User) *FakeUserService {
    return &FakeUserService{fakeUsers: fakeUsers}
}
```

## Integration Testing

TODO...

### Test Helpers

TODO...

### Migrations

Database migrations are managed using the **goose** package.

**Makefile Commands**

- `migrate-create` - Create a new migration file
- `migrate-up` - Apply pending migrations
- `migrate-down` - Rollback last migration

**Migration Execution**

- [ ] Migrations run automatically on application startup (including production)
- [ ] Each deploy triggers auto-migration, no manual DevOps intervention required

**File Naming Convention**

- [ ] Always prefix migration file names with the module name
- [ ] Use snake_case for file names

```bash
# Good - prefixed with module name
auth_init_schema
auth_add_user_roles
platform_init_taskmill

# Bad - missing module prefix
init_schema
add_user_roles
```

**File Organization**

- [ ] All migrations go into a single `./migrations` folder
- [ ] Do NOT separate migrations into subfolders by module

**Required Environment Variables**

The following environment variables must be set for Makefile migration commands:

```bash
POSTGRES_HOST
POSTGRES_PORT
POSTGRES_USER
POSTGRES_PASSWORD
POSTGRES_DB
POSTGRES_SSL
```

These same variables should be used in application configuration for consistency.
