// Package coreserver provides a constructor function to create a core HTTP server with standard middlewares.
package coreserver

import (
	"github.com/rise-and-shine/pkg/alert"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/http/server/middleware"
	"github.com/rise-and-shine/pkg/logger"
)

func New(
	cfg server.Config,
	serviceName, serviceVersion string,
	logger logger.Logger,
	alertPr alert.Provider,
) *server.HTTPServer {
	middlewares := []server.Middleware{
		middleware.NewRecoveryMW(logger),
		middleware.NewTracingMW(),
		middleware.NewTimeoutMW(cfg.HandleTimeout),
		middleware.NewMetaInjectMW(serviceName, serviceVersion),
		middleware.NewAlertingMW(logger, alertPr),
		middleware.NewLoggerMW(logger),
		middleware.NewErrorHandlerMW(cfg.Debug),
	}

	return server.NewHTTPServer(cfg, middlewares)
}
