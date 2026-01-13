package auth

import (
	"go-enterprise-blueprint/internal/modules/auth/ctrl/http"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/internal/portal/auth"

	"github.com/code19m/errx"
	"github.com/rise-and-shine/pkg/http/server"
	"github.com/uptrace/bun"
	"golang.org/x/sync/errgroup"
)

type Config struct {
	Name    string `yaml:"name"    validate:"required"`
	Version string `yaml:"version" validate:"required"`

	HttpServer server.Config `yaml:"http_server" validate:"required"`
}

type Module struct {
	cfg Config

	portal auth.Portal

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
	return m.portal
}

func (m *Module) Start() error {
	var g errgroup.Group

	g.Go(m.httpController.Server().Start)

	err := g.Wait()
	return errx.Wrap(err)
}

func (m *Module) Shutdown() error {
	return errx.Wrap(
		m.httpController.Server().Stop(),
	)
}
