package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/observability/logger"
)

func RunAuth() error {
	app := newApp()
	defer app.shutdown()

	err := app.initAll()
	if err != nil {
		return errx.Wrap(err)
	}

	errChan := make(chan error, 1)
	go func() {
		err := app.auth.Start()
		if err != nil {
			logger.Errorf("error occurred at module.Start: %v", err)
		}
		errChan <- err
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
