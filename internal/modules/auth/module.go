package auth

import (
	"errors"
	"go-enterprise-blueprint/internal/modules/auth/ctrl/asynctask"
	"go-enterprise-blueprint/internal/modules/auth/ctrl/cli"
	"go-enterprise-blueprint/internal/modules/auth/ctrl/consumer"
	"go-enterprise-blueprint/internal/modules/auth/ctrl/http"
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/infra/postgres"
	authportal "go-enterprise-blueprint/internal/modules/auth/portal"
	"go-enterprise-blueprint/internal/modules/auth/usecase"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/createsuperadmin"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/kafka"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Consumers consumer.Config `yaml:"consumers"`
}

type Module struct {
	asynctaskCTRL *asynctask.Controller
	consumerCTRL  *consumer.Controller
	cliCTRL       *cli.Controller
	httpCTRL      *http.Controller

	portal auth.Portal
}

func (m *Module) name() string {
	return "auth"
}

func New(
	cfg Config,
	brokerConfig kafka.BrokerConfig,
	dbConn *bun.DB,
	portalContainer *portal.Container,
	httpServer *server.HTTPServer,
) (*Module, error) {
	var (
		err error
		m   = &Module{}
	)

	// Init repositories
	domainContainer := domain.NewContainer(
		postgres.NewAdminRepo(dbConn),
		postgres.NewSessionRepo(dbConn),
		postgres.NewRoleRepo(dbConn),
		postgres.NewRolePermissionRepo(dbConn),
		postgres.NewActorRoleRepo(dbConn),
		postgres.NewActorPermissionRepo(dbConn),
		postgres.NewUOWFactory(dbConn),
	)

	// Init use cases
	usecaseContainer := usecase.NewContainer(
		createsuperadmin.New(domainContainer),
	)

	// Init portal
	m.portal = authportal.New()

	// Init controllers
	m.cliCTRL = cli.NewController(usecaseContainer)
	m.httpCTRL = http.NewContoller(usecaseContainer, portalContainer, httpServer)
	m.asynctaskCTRL, err = asynctask.NewController(dbConn, m.name(), usecaseContainer)
	if err != nil {
		return nil, errx.Wrap(err)
	}
	m.consumerCTRL, err = consumer.NewController(cfg.Consumers, brokerConfig, usecaseContainer)
	if err != nil {
		return nil, errx.Wrap(err)
	}

	return m, nil
}

func (m *Module) Portal() auth.Portal {
	return m.portal
}

func (m *Module) Start() error {
	var g errgroup.Group

	g.Go(m.asynctaskCTRL.Start)

	g.Go(m.consumerCTRL.Start)

	return errx.Wrap(g.Wait())
}

func (m *Module) Shutdown() error {
	errs := make(chan error, 2) // buffer size == controller count

	go func() { errs <- m.asynctaskCTRL.Shutdown() }()

	go func() { errs <- m.consumerCTRL.Shutdown() }()

	return errx.Wrap(errors.Join(<-errs, <-errs)) // <-errs count == controller count
}

// --- CLI commands of auth module ---

func (m *Module) CreateSuperadmin() error {
	return errx.Wrap(m.cliCTRL.CreateSuperadminCmd())
}
