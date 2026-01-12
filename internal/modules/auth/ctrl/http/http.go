package http

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase"
	"go-enterprise-blueprint/pkg/baseserver"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
)

type Controller struct {
	httpServer *server.HTTPServer
}

func NewContoller(
	cfg server.Config,
	uc *usecase.Container,
) *Controller {
	httpServer := baseserver.New(cfg)

	ctrl := &Controller{httpServer}

	ctrl.httpServer.RegisterRouter(ctrl.initRoutes)

	return ctrl
}

func (c *Controller) Server() *server.HTTPServer {
	return c.httpServer
}

func (c *Controller) initRoutes(r fiber.Router) {
	v1 := r.Group("/auth/v1")

	v1.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{"status": "OK"})
	})

	// Add new handlers here...
}
