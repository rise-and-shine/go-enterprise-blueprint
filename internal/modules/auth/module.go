package auth

import (
	"go-enterprise-blueprint/internal/modules/auth/ctrl/http"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
)

type Config struct {
	Name    string `yaml:"name"    validate:"required"`
	Version string `yaml:"version" validate:"required"`

	HttpServer server.Config `yaml:"http_server" validate:"required"`
}

type Module struct {
	cfg Config

	httpController *http.Controller
}

func New(
	cfg Config,
	dbConn *bun.DB,
	pc *portal.Container,
) *Module {
	m := &Module{cfg: cfg}

	// Init repositories

	// Init pblc

	// Init use cases

	// Init controllers

	return m
}

func (m *Module) Portal() auth.Portal {
	// TODO: return auth portal
	return nil
}

func (m *Module) Start() error {
	return errx.Wrap(
		m.httpController.Server().Start(),
	)
}

func (m *Module) Shutdown() error {
	return errx.Wrap(
		m.httpController.Server().Stop(),
	)
}
