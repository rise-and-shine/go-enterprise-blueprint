package usecase

import "go-enterprise-blueprint/internal/modules/auth/usecase/user/create_superadmin"

type Container struct {
	createSuperAdmin create_superadmin.UseCase
}

func NewContainer(
	createSuperAdmin create_superadmin.UseCase,
) *Container {
	return &Container{
		createSuperAdmin,
	}
}
