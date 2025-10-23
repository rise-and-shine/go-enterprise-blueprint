package app

import (
	"errors"
	"go-enterprise-blueprint/config"
	"log/slog"
	"os"
	"time"

	"go-enterprise-blueprint/pkg/cfgloader"
	"go-enterprise-blueprint/pkg/logger"
)

type App struct{}

type logEntry struct {
	Title   string            `json:"title"`
	Fields  map[string]string `json:"fields,omitempty"`
	Details string            `json:"details,omitempty"`
}

func Build() *App {
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
	}).Info("test info log entry", "entry")

	logger.Warn("i'm just warning you, there is something not ok")

	logger.With("error", errors.New("simulated error")).Error("something went wrong here")

	return &App{}
}

func (a *App) RunAllInOne() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		// Simulate doing work
		slog.Info("running...")
	}
}
