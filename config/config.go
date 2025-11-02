package config

import (
	"go-enterprise-blueprint/config/modconf"

	"github.com/rise-and-shine/pkg/logger"
)

type Config struct {
	Logger logger.Config `yaml:"logger" validate:"required"`

	SecretField string `yaml:"secret_field" secret:"true"`

	// Configurations by modules
	Auth modconf.Auth `yaml:"auth" validate:"required"`
}
