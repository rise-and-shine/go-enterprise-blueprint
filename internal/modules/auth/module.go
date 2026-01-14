package auth

import (
	"go-enterprise-blueprint/internal/modules/auth/ctrl/cli"
	"go-enterprise-blueprint/internal/modules/auth/ctrl/http"
	"go-enterprise-blueprint/internal/modules/auth/domain"
	"go-enterprise-blueprint/internal/modules/auth/infra/postgres"
	"go-enterprise-blueprint/internal/modules/auth/service"
	"go-enterprise-blueprint/internal/modules/auth/service/hasher"
	"go-enterprise-blueprint/internal/modules/auth/service/jwt"
	"go-enterprise-blueprint/internal/modules/auth/usecase"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/admin_login"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/admin_logout"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/admin_refresh_token"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/create_admin"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/disable_admin"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/get_admins"
	"go-enterprise-blueprint/internal/modules/auth/usecase/admin/update_admin"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/get_actor_permissions"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/get_actor_roles"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/get_role_permissions"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/set_actor_permission"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/set_actor_role"
	"go-enterprise-blueprint/internal/modules/auth/usecase/rbac/set_role_permission"
	"go-enterprise-blueprint/internal/modules/auth/usecase/role/create_role"
	"go-enterprise-blueprint/internal/modules/auth/usecase/role/delete_role"
	"go-enterprise-blueprint/internal/modules/auth/usecase/role/get_roles"
	"go-enterprise-blueprint/internal/modules/auth/usecase/role/update_role"
	"go-enterprise-blueprint/internal/modules/auth/usecase/session/delete_user_all_sessions"
	"go-enterprise-blueprint/internal/modules/auth/usecase/session/delete_user_session"
	"go-enterprise-blueprint/internal/modules/auth/usecase/user/create_superadmin"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Name    string `yaml:"name"    validate:"required"`
	Version string `yaml:"version" validate:"required"`

	HttpServer server.Config `yaml:"http_server" validate:"required"`
}

type Module struct {
	cfg Config

	portal auth.Portal

	httpController *http.Controller
	cliController  *cli.Controller
}

func New(
	cfg Config,
	dbConn bun.IDB,
	portalContainer *portal.Container,
) *Module {
	m := &Module{cfg: cfg}

	// Init repositories
	adminRepo := postgres.NewAdminRepo(dbConn)
	sessionRepo := postgres.NewSessionRepo(dbConn)
	roleRepo := postgres.NewRoleRepo(dbConn)
	rolePermissionRepo := postgres.NewRolePermissionRepo(dbConn)
	actorRoleRepo := postgres.NewActorRoleRepo(dbConn)
	actorPermissionRepo := postgres.NewActorPermissionRepo(dbConn)

	dc := domain.NewContainer(
		adminRepo,
		sessionRepo,
		roleRepo,
		rolePermissionRepo,
		actorRoleRepo,
		actorPermissionRepo,
	)

	// Init services
	hasherSvc := hasher.New(0)
	jwtSvc := jwt.New(cfg.JWT)

	sc := service.NewContainer(hasherSvc, jwtSvc)

	// Init use cases
	uc := usecase.NewContainer(
		create_superadmin.New(dc, sc),
		admin_login.New(dc, sc),
		admin_refresh_token.New(dc, sc),
		admin_logout.New(dc),
		delete_user_session.New(dc),
		delete_user_all_sessions.New(dc),
		create_role.New(dc),
		update_role.New(dc),
		delete_role.New(dc),
		get_roles.New(dc),
		set_role_permission.New(dc),
		get_role_permissions.New(dc),
		set_actor_role.New(dc),
		get_actor_roles.New(dc),
		set_actor_permission.New(dc),
		get_actor_permissions.New(dc),
		create_admin.New(dc, sc),
		update_admin.New(dc, sc),
		disable_admin.New(dc),
		get_admins.New(dc),
	)

	// Init controllers
	m.httpController = http.NewContoller(cfg.HttpServer, uc, jwtSvc)
	m.cliController = cli.NewController(uc)

	return m
}

func (m *Module) Portal() auth.Portal {
	return m.portal
}

func (m *Module) CLICommands() []*cobra.Command {
	return m.cliController.Commands()
}

func (m *Module) Start() error {
	var g errgroup.Group

	g.Go(m.httpController.Server().Start)

	err := g.Wait()
	return errx.Wrap(err)
}

func (m *Module) Shutdown() error {
	return errx.Wrap(
		m.httpController.Server().Stop(),
	)
}
