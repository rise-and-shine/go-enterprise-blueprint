package portal

import (
	"go-enterprise-blueprint/internal/portal/auth"
	"go-enterprise-blueprint/internal/portal/esign"
)

// Container holds every modules portal interface.
// It acts as a dependency injection container for the portal layer.
type Container struct {
	auth  auth.Portal
	esign esign.Portal
}

// NewContainer creates a new Container.
func NewContainer(
	auth auth.Portal,
	esign esign.Portal,
) *Container {
	return &Container{
		auth:  auth,
		esign: esign,
	}
}

func (c *Container) Auth() auth.Portal {
	return c.auth
}

func (c *Container) Esign() esign.Portal {
	return c.esign
}
