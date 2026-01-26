package app

import (
	"sync"
	"time"

	"github.com/rise-and-shine/pkg/observability/logger"
)

type shutdownItem struct {
	name string
	fn   func() error
}

// shutdown shuts down high level components first, then infrastructure components.
func (a *app) shutdown() {
	a.shutdownHighLevelComponents()
	a.shutdownInfraComponents()
}

func (a *app) shutdownHighLevelComponents() {
	var items []shutdownItem

	if a.httpServer != nil {
		items = append(items, shutdownItem{name: "http server", fn: a.httpServer.Stop})
	}
	if a.auth != nil {
		items = append(items, shutdownItem{name: "auth module", fn: a.auth.Shutdown})
	}
	// Add your new high level components here...

	if len(items) > 0 {
		a.runShutdown("shutdown_high_level_components", items)
	}
}

func (a *app) shutdownInfraComponents() {
	var items []shutdownItem

	if a.tracerShutdownFunc != nil {
		items = append(items, shutdownItem{name: "trace provider", fn: a.tracerShutdownFunc})
	}
	if a.alertShutdownFunc != nil {
		items = append(items, shutdownItem{name: "alert provider", fn: a.alertShutdownFunc})
	}
	if a.dbConn != nil {
		items = append(items, shutdownItem{name: "database connection", fn: a.dbConn.Close})
	}
	// Add your new infra components here...

	if len(items) > 0 {
		a.runShutdown("shutdown_infra_components", items)
	}
}

func (a *app) runShutdown(operation string, items []shutdownItem) {
	var wg sync.WaitGroup

	for _, item := range items {
		name := item.name
		fn := item.fn

		wg.Go(func() {
			start := time.Now()
			err := fn()
			if err != nil {
				logger.With(
					"component", name,
					"operation", operation,
					"duration", time.Since(start),
				).Errorx(err)
			} else {
				logger.With(
					"component", name,
					"operation", operation,
					"duration", time.Since(start),
				).Info("")
			}
		})
	}

	wg.Wait()
}
