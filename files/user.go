package domain

// import "context"

// // User represents a user entity
// type User struct {
// 	ID    string
// 	Email string
// 	Name  string
// }

// // UserRepository defines persistence operations for users.
// // Mocks will be generated for this interface via .mockery.yaml
// type UserRepository interface {
// 	// GetByID retrieves a user by their unique identifier
// 	GetByID(ctx context.Context, id string) (*User, error)

// 	// GetByEmail retrieves a user by email address
// 	GetByEmail(ctx context.Context, email string) (*User, error)

// 	// Save persists a user (create or update)
// 	Save(ctx context.Context, user *User) error

// 	// Delete removes a user by ID
// 	Delete(ctx context.Context, id string) error

// 	// List returns users with pagination
// 	List(ctx context.Context, offset, limit int) ([]*User, error)
// }

// // Generic repository example (mockery handles generics)
// type Repository[T any] interface {
// 	Get(ctx context.Context, id string) (T, error)
// 	Save(ctx context.Context, entity T) error
// 	Delete(ctx context.Context, id string) error
// }
