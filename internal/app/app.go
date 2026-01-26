package app

import (
	"go-enterprise-blueprint/internal/modules/auth"

	"github.com/rise-and-shine/pkg/cfgloader"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/rise-and-shine/pkg/kafka"
	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/observability/tracing"
	"github.com/rise-and-shine/pkg/pg"
	"github.com/uptrace/bun"
)

type Config struct {
	// --- Shared configs ---

	Service ServiceConfig `yaml:"service" validate:"required"`

	Logger logger.Config `yaml:"logger" validate:"required"`

	Tracing tracing.Config `yaml:"tracing" validate:"required"`

	Alert alert.Config `yaml:"alert" validate:"required"`

	Postgres pg.Config `yaml:"postgres" validate:"required"`

	KafkaBroker kafka.BrokerConfig `yaml:"kafka_broker" validate:"required"`

	HTTPServer server.Config `yaml:"http_server" validate:"required"`

	// --- Module specific configs ---

	Auth auth.Config `yaml:"auth"`
}

type app struct {
	cfg Config

	dbConn             *bun.DB
	tracerShutdownFunc func() error
	alertShutdownFunc  func() error

	httpServer *server.HTTPServer

	auth *auth.Module
}

func newApp() *app {
	app := &app{
		cfg: cfgloader.MustLoad[Config](),
	}
	return app
}

type ServiceConfig struct {
	// Name is the name of the service
	Name string `json:"name" validate:"required"`

	// Version is the version of the service
	Version string `json:"version" validate:"required"`
}
