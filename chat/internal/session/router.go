package session

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/middleware"
)

type sessionHandler interface {
	NewSession(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Rename(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

func SetupSessionRoutes(app *fiber.App, h sessionHandler, kc *gocloak.GoCloak, realm, cipher string) {
	g := app.Group("/session")
	g.Use(middleware.JWT(kc, realm, cipher))

	g.Post("/new", h.NewSession)
	g.Get("/list", h.List)
	g.Put("/rename", h.Rename)
	g.Delete("/:id", h.Delete)
}
