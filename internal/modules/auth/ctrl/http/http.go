package http

import (
	"go-enterprise-blueprint/internal/modules/auth/usecase"
	"go-enterprise-blueprint/internal/portal"
	"go-enterprise-blueprint/pkg/baseserver"

	"github.com/gofiber/fiber/v2"
	"github.com/rise-and-shine/pkg/http/server"
)

type Controller struct {
	usecaseContainer *usecase.Container
	portalContainer  *portal.Container
	httpServer       *server.HTTPServer
}

func NewContoller(
	serverConfig server.Config,
	usecaseContainer *usecase.Container,
	portalContainer *portal.Container,
) *Controller {
	ctrl := &Controller{
		usecaseContainer,
		portalContainer,
		baseserver.New(serverConfig),
	}

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

	// Add your routes here...
}
