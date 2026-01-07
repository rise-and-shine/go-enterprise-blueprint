package http

import (
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/service"
	"go-enterprise-blueprint/internal/modules/auth/usecase/user/create_user"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/pkg/baseserver"

	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/http/server/forward"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
)

type Server struct {
	core *server.HTTPServer
}

func New(
	cfg server.Config,
	serviceName, serviceVersion string,
	logger logger.Logger,
	alertPr alert.Provider,
) *Server {
	s := &Server{
		baseserver.New(
			cfg,
			serviceName,
			serviceVersion,
			logger,
			alertPr,
		),
	}
	s.initRoutes()
	return s
}

func (s *Server) Start() error {
	return s.core.Start()
}

func (s *Server) Stop() error {
	return s.core.Stop()
}

func (s *Server) initRoutes() {
	app := s.core.GetApp()

	app.Post("/", forward.ToUseCase(create_user.New(domain.Container{}, service.Container{}, portal.Container{})))
}
