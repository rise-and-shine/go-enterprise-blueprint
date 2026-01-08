package modules

import "github.com/rise-and-shine/pkg/http/server"

type Auth struct {
	HttpServer server.Config `yaml:"http_server" validate:"required"`
}
