package main

import (
	"github.com/spf13/cobra"
)

func main() {
	var root = &cobra.Command{}

	// all in one
	root.AddCommand(runAllInOne())

	// reverse proxy
	// root.AddCommand(runProxy())

	// ...add http server commands here...

	// cron manager
	// ...add cron manager run command here...

	// module based cli commands
	// ...add cli commands here...

	// ignoring error since it's already displayed by cobra.
	_ = root.Execute()
}

func runAllInOne() *cobra.Command {
	return &cobra.Command{
		Use:   "run-all-in-one",
		Short: "Run all services in one process",
		Run: func(_ *cobra.Command, _ []string) {
			// app.RunAllInOne()
		},
	}
}

// func runProxy() *cobra.Command {
// 	var configPath string

// 	cmd := &cobra.Command{
// 		Use:   "run-proxy",
// 		Short: "Run the reverse proxy server",
// 		Run: func(_ *cobra.Command, _ []string) {
// 			logger.SetGlobal(logger.Config{
// 				Level:    "info",
// 				Encoding: "pretty",
// 			})

// 			cfg, err := easyproxy.LoadConfig(configPath)
// 			if err != nil {
// 				log.Fatalf("Failed to load config: %v", err)
// 			}

// 			proxy, err := easyproxy.New(cfg)
// 			if err != nil {
// 				log.Fatalf("Failed to create proxy: %v", err)
// 			}

// 			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
// 			defer cancel()

// 			if err := proxy.Run(ctx); err != nil {
// 				log.Fatalf("Proxy server error: %v", err)
// 			}
// 		},
// 	}

// 	cmd.Flags().StringVarP(&configPath, "config", "c", "proxy.yaml", "Path to proxy config file")

// 	return cmd
// }
