package app

import (
	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/observability/tracing"
)

type app struct {
	cfg Config

	tracerShutdownFunc func() error
	alertShutdownFunc  func() error
	dbCloser           func() error
}

func newAppWithConfig(cfg Config) *app {
	return &app{
		cfg: cfg,
	}
}

func (a *app) initObservability(
	serviceName string,
	serviceVersin string,
) error {
	var err error

	// set global service information
	meta.SetServiceInfo(serviceName, serviceVersin)

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

	return nil
}

func (a *app) shutdown() {
	// TODO...
}
