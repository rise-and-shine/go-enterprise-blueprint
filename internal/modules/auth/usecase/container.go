package usecase

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/createsuperadmin"
)

type Container struct {
	createSuperadmin createsuperadmin.UseCase
}

func NewContainer(
	createSuperadmin createsuperadmin.UseCase,
) *Container {
	return &Container{
		createSuperadmin: createSuperadmin,
	}
}

func (c *Container) CreateSuperadmin() createsuperadmin.UseCase {
	return c.createSuperadmin
}
