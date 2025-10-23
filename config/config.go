package config

import "go-enterprise-blueprint/pkg/logger"

type Config struct {
	Logger logger.Config `yaml:"logger" validate:"required"`

	SecretField string `yaml:"secret_field" secret:"true"`
}
