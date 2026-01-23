package main

import (
	"go-enterprise-blueprint/internal/app"

	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/spf13/cobra"
)

func main() {
	var root = &cobra.Command{}

	root.AddCommand(run())

	root.AddCommand(app.AuthCommands())
	// Add new modules CLI commands here...

	// ignoring error since it's already displayed by cobra.
	_ = root.Execute()
}

// run registers a main command that runs all services.
func run() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run all services in one process",
		Run: func(_ *cobra.Command, _ []string) {
			err := app.Run()
			if err != nil {
				logger.Fatalx(err)
			}
		},
	}
}
