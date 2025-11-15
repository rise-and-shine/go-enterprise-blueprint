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

// UseCaseRegistry provides a convenient way to register usecases without repetitive code.
type UseCaseRegistry struct {
	api             huma.API
	registrationFns []func()
}

// NewUseCaseRegistry creates a new usecase registry.
func NewUseCaseRegistry(api huma.API) *UseCaseRegistry {
	return &UseCaseRegistry{
		api:             api,
		registrationFns: make([]func(), 0),
	}
}

// registerOptions holds configuration options for registering usecases.
type registerOptions struct {
	Security    openAPISecurity
	Method      string
	Path        string
	OperationID string
}

// registerOption is a function that modifies RegisterOptions.
type registerOption func(*registerOptions)

// WithSecurity sets the security configuration.
func WithSecurity(sec openAPISecurity) registerOption {
	return func(opts *registerOptions) {
		opts.Security = sec
	}
}

// WithPath sets the HTTP path.
func WithPath(path string) registerOption {
	return func(opts *registerOptions) {
		opts.Path = path
	}
}

// AddWrite registers a write usecase with the given API and options.
func AddWrite[I, O, S any](api huma.API, uc ucdef.UserWriteAction[I, O, S], opts ...registerOption) {
	// Set default options
	options := registerOptions{
		Security:    BearerSecurity(),
		Method:      http.MethodPost,
		Path:        uc.OperationID(),
		OperationID: uc.OperationID(),
	}

	// Apply provided options
	for _, opt := range opts {
		opt(&options)
	}

	huma.Register(api,
		huma.Operation{
			OperationID: options.OperationID,
			Method:      options.Method,
			Path:        options.Path,
			Security:    options.Security,
		},
		func(ctx context.Context, i *I) (*O, error) {
			// Placeholder implementation
			state, err := uc.Validate(ctx, i)
			if err != nil {
				return nil, errx.Wrap(err)
			}

			auditInfo := uc.AuditInfo(state)
			logger.Named("auditor").WithContext(ctx).With("audit_info", auditInfo).Info()

			out, err := uc.Execute(ctx, state)
			return out, errx.Wrap(err)
		},
	)
}

// AddRead registers a read usecase with the given API and options.
func AddRead[I, O, S any](api huma.API, uc ucdef.UserReadAction[I, O, S], opts ...registerOption) {
	// Set default options
	options := registerOptions{
		Security:    BearerSecurity(),
		Method:      http.MethodGet,
		Path:        uc.OperationID(),
		OperationID: uc.OperationID(),
	}

	// Apply provided options
	for _, opt := range opts {
		opt(&options)
	}

	huma.Register(api,
		huma.Operation{
			OperationID: options.OperationID,
			Method:      options.Method,
			Path:        options.Path,
			Security:    options.Security,
			Middlewares: huma.Middlewares{
				func(ctx huma.Context, next func(huma.Context)) {
					// before
					next(ctx)
					// after
				},
			},
		},
		func(ctx context.Context, i *I) (*O, error) {
			// Placeholder implementation
			state, err := uc.Validate(ctx, i)
			if err != nil {
				return nil, errx.Wrap(err)
			}

			// Note: UserReadAction doesn't have AuditInfo method
			out, err := uc.Execute(ctx, state)
			return out, errx.Wrap(err)
		},
	)
}

func (s *Server) InitRoutes() {
	// Here starts Huma duties. huma.rocks

	app := s.core.GetApp()
	api := humafiber.New(app, huma.DefaultConfig("Go Enterprise Blueprint Auth Service", "v1.0.0"))

	// Add write usecases with the new API
	uc1001 := uc1001createuser.New(domain.Container{}, service.Container{}, portal.Container{})
	AddWrite(api, uc1001) // Uses default options: POST, BearerSecurity, operationID from usecase

	// Example with custom options:
	// AddWrite(api, uc1001,
	//     WithSecurity(NoSecurity()),
	//     WithPath("/custom/path"),
	//     WithMethod(http.MethodPut),
	// )

	// Add more usecases here:
	// uc1002 := uc1002updateuser.New(domain.Container{}, service.Container{}, portal.Container{})
	// AddWrite(api, uc1002, WithSecurity(BearerSecurity()))

	// uc2001 := uc2001getuser.New(domain.Container{}, service.Container{}, portal.Container{})
	// AddRead(api, uc2001) // Uses default options: GET, BearerSecurity, operationID from usecase
}

type AuthMethod interface {
	OpenAPISecurity() []map[string][]string
	AuthMiddleware() []func(ctx huma.Context, next func(huma.Context))
}

const (
	Bearer AuthMethod = "bearer"
	NoAuth AuthMethod = "no_auth"
)

type openAPISecurity []map[string][]string

func BearerSecurity() openAPISecurity {
	return openAPISecurity{{"bearer": {}}}
}

func NoSecurity() openAPISecurity {
	return openAPISecurity{}
}
