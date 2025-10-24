package main

import (
	"go-enterprise-blueprint/internal/app"

	"github.com/spf13/cobra"
)

func main() {
	var root = &cobra.Command{}

	// all in one
	root.AddCommand(allInOne())

	// http servers
	// ...add http server commands here...

	// cli commands
	// ...add cli commands here...

	// cron manager
	// ...add cron manager run command here...

	// ignoring error since it's already displayed by cobra.
	_ = root.Execute()
}

func allInOne() *cobra.Command {
	return &cobra.Command{
		Use:   "all-in-one",
		Short: "Run all services in one process",
		Run: func(_ *cobra.Command, _ []string) {
			app.Build().RunAllInOne()
		},
	}
}
