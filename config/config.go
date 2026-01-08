package config

import (
	"go-enterprise-blueprint/config/modules"

	"github.com/rise-and-shine/pkg/observability/logger"
)

type Config struct {
	// --- Shared configs ---

	Logger logger.Config `yaml:"logger" validate:"required"`

	SecretField string `yaml:"secret_field" secret:"true"`

	// --- Module specific configs ---

	Auth modules.Auth `yaml:"auth" validate:"required"`
}
