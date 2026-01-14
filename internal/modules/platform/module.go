package platform

import (
	"time"

	"go-enterprise-blueprint/internal/modules/platform/usecase/get_docs"
	"go-enterprise-blueprint/pkg/baseserver"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
)

type Config struct {
	Name        string        `yaml:"name"        validate:"required"`
	Version     string        `yaml:"version"     validate:"required"`
	Description string        `yaml:"description"`
	HttpServer  server.Config `yaml:"http_server" validate:"required"`
}

type Module struct {
	cfg        Config
	httpServer *server.HTTPServer
	getDocsUC  get_docs.UseCase
	startTime  time.Time
}

func New(cfg Config) *Module {
	m := &Module{
		cfg:       cfg,
		startTime: time.Now(),
	}

	m.getDocsUC = get_docs.New(get_docs.Config{
		ServiceName:        cfg.Name,
		ServiceVersion:     cfg.Version,
		ServiceDescription: cfg.Description,
		StartTime:          m.startTime,
	})

	m.httpServer = baseserver.New(cfg.HttpServer)
	m.httpServer.RegisterRouter(m.initRoutes)

	return m
}

func (m *Module) initRoutes(r fiber.Router) {
	r.Get("/platform/get-docs", m.getDocs)
}

func (m *Module) getDocs(ctx *fiber.Ctx) error {
	output, err := m.getDocsUC.Execute(ctx.Context(), &get_docs.Input{})
	if err != nil {
		return err
	}
	return ctx.JSON(output)
}

func (m *Module) Start() error {
	return m.httpServer.Start()
}

func (m *Module) Shutdown() error {
	return m.httpServer.Stop()
}
