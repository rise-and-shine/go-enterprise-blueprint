package app

import (
	"time"

	"github.com/rise-and-shine/pkg/observability/logger"
)

func RunAllInOne(cfg Config) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Simulate doing work
		logger.Info("running...")
	}

	return nil
}
