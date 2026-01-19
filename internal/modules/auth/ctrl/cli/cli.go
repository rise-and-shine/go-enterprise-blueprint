// Package cli provides container of cobra CLI commands for auth module.
package cli

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase"
)

type Controller struct {
	usecaseContainer *usecase.Container
}

func NewController(usecaseContainer *usecase.Container) *Controller {
	return &Controller{
		usecaseContainer,
	}
}
