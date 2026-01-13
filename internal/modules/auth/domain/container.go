package domain

// Container holds domain interfaces.
// It acts as a dependency injection container for the domain layer.
type Container struct {
}

func NewContainer() *Container {
	return &Container{}
}
