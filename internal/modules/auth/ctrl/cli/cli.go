// Package cli provides container of cobra CLI commands for auth module.
package cli

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase"

	"github.com/spf13/cobra"
)

type Controller struct {
	usecaseContainer *usecase.Container
}

func NewController(usecaseContainer *usecase.Container) *Controller {
	return &Controller{
		usecaseContainer,
	}
}

func (c *Controller) Commands() []*cobra.Command {
	return []*cobra.Command{
		c.createSuperadminCmd(),
	}
}
