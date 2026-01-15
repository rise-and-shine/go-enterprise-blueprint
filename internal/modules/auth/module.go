package auth

import (
	"go-enterprise-blueprint/internal/modules/auth/ctrl/asynctask"
	"go-enterprise-blueprint/internal/modules/auth/ctrl/cli"
	"go-enterprise-blueprint/internal/modules/auth/ctrl/consumer"
	"go-enterprise-blueprint/internal/modules/auth/ctrl/http"
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/infra/postgres"
	authportal "go-enterprise-blueprint/internal/modules/auth/portal"
	"go-enterprise-blueprint/internal/modules/auth/usecase"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/create_superadmin"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/kafka"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Name    string `yaml:"name"    validate:"required"`
	Version string `yaml:"version" validate:"required"`

	HttpServer server.Config   `yaml:"http_server" validate:"required"`
	Consumers  consumer.Config `yaml:"consumers"   validate:"required"`
}

type Module struct {
	asynctaskCTRL *asynctask.Controller
	consumerCTRL  *consumer.Controller
	cliCTRL       *cli.Controller
	httpCTRL      *http.Controller

	portal auth.Portal
}

func New(
	cfg Config,
	brokerConfig kafka.BrokerConfig,
	dbConn *bun.DB,
	portalContainer *portal.Container,
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
	)

	// Init use cases
	usecaseContainer := usecase.NewContainer(
		create_superadmin.New(domainContainer),
	)

	// Init portal
	m.portal = authportal.New()

	// Init controllers
	m.cliCTRL = cli.NewController(usecaseContainer)
	m.httpCTRL = http.NewContoller(cfg.HttpServer, usecaseContainer, portalContainer)
	m.asynctaskCTRL, err = asynctask.NewController(dbConn, cfg.Name, usecaseContainer)
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

func (m *Module) CLICommands() []*cobra.Command {
	return m.cliCTRL.Commands()
}

func (m *Module) Start() error {
	var g errgroup.Group

	g.Go(m.asynctaskCTRL.Start)
	g.Go(m.consumerCTRL.Start)
	g.Go(m.httpCTRL.Server().Start)

	err := g.Wait()
	return errx.Wrap(err)
}

func (m *Module) Shutdown() error {
	return errx.Wrap(
		m.httpCTRL.Server().Stop(),
	)
}
