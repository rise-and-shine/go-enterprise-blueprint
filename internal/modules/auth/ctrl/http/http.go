package http

import (
	"context"
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/service"
	"go-enterprise-blueprint/internal/modules/auth/usecase/user/uc1001createuser"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/pkg/coreserver"
	"go-enterprise-blueprint/pkg/ucdef"
	"net/http"

	"github.com/code19m/errx"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humafiber"
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
	core := coreserver.New(
		cfg,
		serviceName,
		serviceVersion,
		logger,
		alertPr,
	)
	return &Server{
		core: core,
	}
}

func (s *Server) Start() error {
	return s.core.Start()
}

func (s *Server) Stop() error {
	return s.core.Stop()
}

func (s *Server) InitRoutes() {
	// Here starts Huma duties. huma.rocks

	app := s.core.GetApp()
	api := humafiber.New(app, huma.DefaultConfig("Go Enterprise Blueprint Auth Service", "v1.0.0"))

	uc1001 := uc1001createuser.New(domain.Container{}, service.Container{}, portal.Container{})

	ForwardWrite(api, NewWriteHandler(uc1001, BearerSecurity()))
}

func ForwardWrite[I, O, S any](api huma.API, handler WriteHandler[I, O, S]) {
	huma.Register(api,
		huma.Operation{
			OperationID: handler.UseCase.OperationID(),
			Method:      http.MethodPost,
			Path:        handler.UseCase.OperationID(),
		},
		Executor(handler.UseCase),
	)
}

func Executor[I, O, S any](uc ucdef.UserWriteAction[I, O, S]) func(ctx context.Context, i *I) (*O, error) {
	return func(ctx context.Context, i *I) (*O, error) {
		state, err := uc.Validate(ctx, i)
		if err != nil {
			return nil, errx.Wrap(err)
		}

		auditInfo := uc.AuditInfo(state)
		// TODO: produce audit info
		logger.Named("auditor").WithContext(ctx).With("audit_info", auditInfo).Info()

		out, err := uc.Execute(ctx, state)

		return out, errx.Wrap(err)
	}
}

type WriteHandler[I, O, S any] struct {
	UseCase ucdef.UserWriteAction[I, O, S]

	Security openAPISecurity // default is bearer
}

func NewWriteHandler[I, O, S any](useCase ucdef.UserWriteAction[I, O, S], sec openAPISecurity) WriteHandler[I, O, S] {
	return WriteHandler[I, O, S]{useCase, sec}
}

type openAPISecurity []map[string][]string

func BearerSecurity() openAPISecurity {
	return openAPISecurity{{"bearer": {}}}
}

func NoSecurity() openAPISecurity {
	return openAPISecurity{}
}
