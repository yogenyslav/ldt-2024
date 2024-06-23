package session

import (
	"github.com/gofiber/fiber/v2"
)

type sessionHandler interface {
	NewSession(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	Rename(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
}

// SetupSessionRoutes маппит пути для сессий.
func SetupSessionRoutes(app fiber.Router, h sessionHandler) {
	g := app.Group("/session")

	g.Post("/new", h.NewSession)
	g.Get("/list", h.List)
	g.Get("/:id", h.FindOne)
	g.Put("/rename", h.Rename)
	g.Delete("/:id", h.Delete)
}
