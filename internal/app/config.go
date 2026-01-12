package app

import (
	"go-enterprise-blueprint/internal/modules/auth"

	"github.com/rise-and-shine/pkg/observability/alert"
	"github.com/rise-and-shine/pkg/observability/logger"
	"github.com/rise-and-shine/pkg/observability/tracing"
	"github.com/rise-and-shine/pkg/pg"
)

type Config struct {
	// --- Shared configs ---

	Logger logger.Config `yaml:"logger" validate:"required"`

	Tracing tracing.Config `yaml:"tracing" validate:"required"`

	Alert alert.Config `yaml:"alert" validate:"required"`

	Postgres pg.Config `yaml:"postgres" validate:"required"`

	// --- Module specific configs ---

	Auth auth.Config `yaml:"auth" validate:"required"`
}
