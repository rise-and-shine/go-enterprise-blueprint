package app

import (
	"errors"
	"go-enterprise-blueprint/config"
	"log/slog"
	"os"
	"time"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/cfgloader"
	"github.com/rise-and-shine/pkg/meta"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/observability/tracing"
)

type app struct {
	cfg config.Config

	tracerShutdownFunc func() error
	alertShutdownFunc  func() error
	dbCloser           func() error
}

func newAppWithConfig(cfg config.Config) *app {
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

func initApp() *app {
	app := &app{}

	cfg := cfgloader.MustLoad[config.Config]()

	logger, err := logger.New(cfg.Logger)
	if err != nil {
		slog.Error("failed to initialize logger", "error", err)
		os.Exit(1)
	}

	logger = logger.Named("application").With("method", "build")

	logger.Debug("start testing logger")

	logger.With("entry", logEntry{
		Title: "Hi",
		Fields: map[string]string{
			"field1": "value1",
			"field2": "value2",
		},
		Details: "This is a detailed message for the info log entry.",
	}).Info("test info log entry")

	logger.Warn("i'm just warning you, there is something not ok")

	logger.With("error", errors.New("simulated error")).Error("something went wrong here")

	err = errx.New(
		"qwer",
		errx.WithCode("qwer_code"),
		errx.WithType(errx.T_Throttling),
	)

	err = errx.Wrap(err)
	err = errx.Wrap(err)

	logger.Warnx(err)

	return &App{}
}

func (a *App) RunAllInOne() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Simulate doing work
		logger.Info("running...")
	}
}
