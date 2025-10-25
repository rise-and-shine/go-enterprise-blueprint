package portal

import (
	"go-enterprise-blueprint/internal/portal/auth"
	"go-enterprise-blueprint/internal/portal/doc"
	"go-enterprise-blueprint/internal/portal/esign"
	"go-enterprise-blueprint/internal/portal/filestore"
)

// Container holds every modules portal interface.
// It acts as a dependency injection container for the portal layer.
type Container struct {
	auth      auth.Portal
	doc       doc.Portal
	esign     esign.Portal
	filestore filestore.Portal
}

// NewContainer creates a new Container.
func NewContainer(
	auth auth.Portal,
	doc doc.Portal,
	esign esign.Portal,
	filestore filestore.Portal,
) *Container {
	return &Container{
		auth:      auth,
		doc:       doc,
		esign:     esign,
		filestore: filestore,
	}
}

func (c *Container) Auth() auth.Portal {
	return c.auth
}

func (c *Container) Doc() doc.Portal {
	return c.doc
}

func (c *Container) Esign() esign.Portal {
	return c.esign
}

func (c *Container) Filestore() filestore.Portal {
	return c.filestore
}
