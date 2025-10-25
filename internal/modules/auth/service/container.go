package service

// Container holds current services interfaces.
// It acts as a dependency injection container for the service layer.
type Container struct{}

// NewContainer creates a new Container.
func NewContainer() *Container {
	return &Container{}
}
