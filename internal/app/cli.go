package app

import (
	"github.com/code19m/errx"
	"github.com/spf13/cobra"
)

func AuthCommands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Auth module CLI commands",
	}

	cmd.AddCommand(createSuperAdminCmd())
	// Add auth modules new CLI commands here...

	return cmd
}

func createSuperAdminCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "create-superadmin",
		Short: "Create superadmin account for system bootstrap",
		RunE: func(_ *cobra.Command, _ []string) error {
			app := newApp()
			defer app.shutdownInfraComponents()

			err := app.init()
			if err != nil {
				return errx.Wrap(err)
			}

			return app.auth.CreateSuperadmin()
		},
	}
}

// Add your new CLI commands here...
