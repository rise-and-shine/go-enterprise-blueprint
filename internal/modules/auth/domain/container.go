package domain

// Container holds domain interfaces.
// It acts as a dependency injection container for the domain layer.
type Container struct {
	userRepo user.Repository

	uowFactory uow.Factory
}

func (c *Container) UserRepo() user.Repository {
	return c.userRepo
}

func (c *Container) UowFactory() uow.Factory {
	return c.uowFactory
}

// NewContainer creates a new Container.
func NewContainer(
	userRepo user.Repository, 
	uowFactory uow.Factory,
) *Container {
	return &Container{
		userRepo,
		uowFactory,
	}
}
