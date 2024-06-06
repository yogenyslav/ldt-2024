package session

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/middleware"
)

type sessionHandler interface {
	NewSession(c *fiber.Ctx) error
}

func SetupSessionRoutes(app *fiber.App, h sessionHandler, kc *gocloak.GoCloak, realm, cipher string) {
	g := app.Group("/session")
	g.Use(middleware.JWT(kc, realm, cipher))

	g.Post("/new", h.NewSession)
}
