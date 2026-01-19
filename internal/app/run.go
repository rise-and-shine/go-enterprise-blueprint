package app

import (
	"go-enterprise-blueprint/internal/modules/auth"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/pkg/baseserver"
	"os"
	"os/signal"
	"syscall"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/observability/tracing"
	"github.com/rise-and-shine/pkg/pg"
	"golang.org/x/sync/errgroup"
)

func Run() error {
	app := newApp()
	defer app.shutdown()

	err := app.init()
	if err != nil {
		return errx.Wrap(err)
	}

	errChan := make(chan error)

	// run all high level components
	go func() {
		errChan <- app.runHighLevelComponents()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	// error occurred at module.Start
	case err := <-errChan:
		return errx.Wrap(err)

	// signal received, just return nil to trigger app.shutdown()
	case <-quit:
		return nil
	}
}

func (a *app) runHighLevelComponents() error {
	var g errgroup.Group

	g.Go(a.httpServer.Start)

	// Run your modules here...
	g.Go(a.auth.Start)

	return errx.Wrap(g.Wait())
}

func (a *app) init() error {
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

	// init http server
	a.httpServer = baseserver.New(a.cfg.HttpServer)

	return nil
}

func (a *app) initModules() error {
	var (
		err error
	)

	portalContainer := &portal.Container{}

	// Init all your modules here...
	a.auth, err = auth.New(
		a.cfg.Auth, a.cfg.KafkaBroker, a.dbConn, portalContainer, a.httpServer,
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
