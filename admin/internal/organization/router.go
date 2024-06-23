package organization

import (
	"github.com/gofiber/fiber/v2"
)

type organizationHandler interface {
	InsertOne(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	ImportData(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
}

// SetupOrganizationRoutes устанавливает маршруты для организаций.
func SetupOrganizationRoutes(app fiber.Router, h organizationHandler) {
	g := app.Group("/organization")

	g.Post("/", h.InsertOne)
	g.Get("/", h.FindOne)
	g.Post("/import", h.ImportData)
	g.Put("/", h.UpdateOne)
}
