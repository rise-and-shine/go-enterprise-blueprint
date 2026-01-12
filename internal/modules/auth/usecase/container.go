package usecase

import "go-enterprise-blueprint/internal/modules/auth/usecase/user/create_user"

type Container struct {
	createSuperUser create_user.UseCase
}
