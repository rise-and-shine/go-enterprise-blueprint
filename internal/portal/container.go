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

func (c *Container) SetAuthPortal(auth auth.Portal) {
	c.auth = auth
}

func (c *Container) SetEsignPortal(esign esign.Portal) {
	c.esign = esign
}

func (c *Container) Auth() auth.Portal {
	return c.auth
}

func (c *Container) Esign() esign.Portal {
	return c.esign
}
