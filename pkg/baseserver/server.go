// Package baseserver provides a constructor function to create a base HTTP server with standard middlewares.
package baseserver

import (
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/http/server/middleware"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
)

func New(
	cfg server.Config,
	serviceName, serviceVersion string,
	logger logger.Logger,
	alertPr alert.Provider,
) *server.HTTPServer {
	middlewares := []server.Middleware{
		middleware.NewRecoveryMW(cfg.Debug),
		middleware.NewTracingMW(),
		middleware.NewTimeoutMW(cfg.HandleTimeout),
		middleware.NewMetaInjectMW(),
		middleware.NewAlertingMW(),
		middleware.NewLoggerMW(cfg.Debug),
		middleware.NewErrorHandlerMW(cfg.Debug),
	}

	return server.NewHTTPServer(cfg, middlewares)
}
