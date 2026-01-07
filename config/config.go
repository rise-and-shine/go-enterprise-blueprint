package config

import (
	"github.com/rise-and-shine/pkg/observability/logger"
)

type Config struct {
	Logger logger.Config `yaml:"logger" validate:"required"`

	SecretField string `yaml:"secret_field" secret:"true"`
}
