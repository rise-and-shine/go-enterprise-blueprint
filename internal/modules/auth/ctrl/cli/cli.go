// Package cli provides container of cobra CLI commands for auth module.
package cli

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase"

	"github.com/spf13/cobra"
)

type Controller struct {
	ucContainer *usecase.Container
}

func NewController(ucContainer *usecase.Container) *Controller {
	return &Controller{
		ucContainer,
	}
}

func (c *Controller) Commands() []*cobra.Command {
	return []*cobra.Command{
		c.createSuperuserCmd(),
	}
}
