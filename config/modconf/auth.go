package modconf

import "github.com/rise-and-shine/pkg/http/server"

type Auth struct {
	HTTPServer server.Config `yaml:"http_server" validate:"required"`
}
