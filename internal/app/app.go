package app

import (
	"go-enterprise-blueprint/internal/modules/auth"
	"go-enterprise-blueprint/internal/portal"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/cfgloader"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/observability/tracing"
	"github.com/rise-and-shine/pkg/pg"
	"github.com/uptrace/bun"
)

type app struct {
	cfg Config

	dbConn             *bun.DB
	tracerShutdownFunc func() error
	alertShutdownFunc  func() error

	auth *auth.Module
}

func newApp() *app {
	app := &app{
		cfg: cfgloader.MustLoad[Config](),
	}
	return app
}

func (a *app) initAll() error {
	err := a.initSharedComponents(a.cfg.Auth.Name, a.cfg.Auth.Version)
	if err != nil {
		return errx.Wrap(err)
	}

	err = a.initModules()
	return errx.Wrap(err)
}

func (a *app) initSharedComponents(
	serviceName string,
	serviceVersion string,
) error {
	var err error

	// set global service information
	meta.SetServiceInfo(serviceName, serviceVersion)

	// init logger
	logger.SetGlobal(a.cfg.Logger)

	// init metrics
	// Metrics provider not implemented yet...

	// init tracing
	a.tracerShutdownFunc, err = tracing.InitGlobalTracer(a.cfg.Tracing)
	if err != nil {
		return errx.Wrap(err)
	}

	// init alerting
	err = alert.SetGlobal(a.cfg.Alert)
	if err != nil {
		return errx.Wrap(err)
	}
	a.alertShutdownFunc = alert.ShutdownGlobal

	// init db connection pool
	a.dbConn, err = pg.NewBunDB(a.cfg.Postgres)
	if err != nil {
		return errx.Wrap(err)
	}

	return nil
}

func (a *app) initModules() error {
	var (
		err error
	)

	portalContainer := &portal.Container{}

	// Init all your modules here...
	a.auth, err = auth.New(
		a.cfg.Auth, a.cfg.KafkaBroker, a.dbConn, portalContainer,
	)
	if err != nil {
		return errx.Wrap(err)
	}

	// Audit

	// Esign

	// Platform

	// Set all portal implementations here...
	portalContainer.SetAuthPortal(a.auth.Portal())
	// portalContainer.SetAuditPortal(audit.Portal())
	// portalContainer.SetEsignPortal(esign.Portal())
	// portalContainer.SetPlatformPortal(platform.Portal())

	return nil
}

func (a *app) shutdown() {
	// modulle

	if a.alertShutdownFunc != nil {
		a.alertShutdownFunc()
	}

	// if a.dbConn !-
}
