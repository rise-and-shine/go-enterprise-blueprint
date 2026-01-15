package usecase

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/create_superadmin"
)

type Container struct {
	createSuperadmin create_superadmin.UseCase
}

func NewContainer(
	createSuperadmin create_superadmin.UseCase,
) *Container {
	return &Container{
		createSuperadmin: createSuperadmin,
	}
}

func (c *Container) CreateSuperadmin() create_superadmin.UseCase {
	return c.createSuperadmin
}
