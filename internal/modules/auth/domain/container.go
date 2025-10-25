package domain

// Container holds current domains interfaces.
// It acts as a dependency injection container for the domain layer.
type Container struct{}

// NewContainer creates a new Container.
func NewContainer() *Container {
	return &Container{}
}
