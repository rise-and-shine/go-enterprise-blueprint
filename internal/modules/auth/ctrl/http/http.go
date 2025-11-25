package http

import (
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/service"
	"go-enterprise-blueprint/internal/modules/auth/usecase/user/uc1001createuser"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/pkg/coreserver"
	"go-enterprise-blueprint/pkg/http/server/forward"

	"github.com/rise-and-shine/pkg/alert"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/logger"
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
		coreserver.New(
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

	app.Post("/", forward.ToUseCase(uc1001createuser.New(domain.Container{}, service.Container{}, portal.Container{})))
}
