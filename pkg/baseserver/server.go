// Package baseserver provides a constructor function to create a base HTTP server with standard middlewares.
package baseserver

import (
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/http/server/middleware"
)

func New(
	cfg server.Config,
) *server.HTTPServer {
	middlewares := []server.Middleware{
		middleware.NewRecoveryMW(cfg.HideErrorDetails),
		middleware.NewTracingMW(),
		middleware.NewTimeoutMW(cfg.HandleTimeout),
		middleware.NewAlertingMW(),
		middleware.NewLoggerMW(cfg.HideErrorDetails),
		middleware.NewErrorHandlerMW(cfg.HideErrorDetails),
	}

	return server.NewHTTPServer(cfg, middlewares)
}
