package user

import (
	"github.com/gofiber/fiber/v2"
)

type userHandler interface {
	NewUser(c *fiber.Ctx) error
	InsertOrganization(c *fiber.Ctx) error
	DeleteOrganization(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
}

// SetupUserRoutes инициализирует роуты для работы с пользователями.
func SetupUserRoutes(app fiber.Router, h userHandler) {
	g := app.Group("/user")

	g.Post("/", h.NewUser)
	g.Post("/organization", h.InsertOrganization)
	g.Delete("/", h.DeleteOrganization)
	g.Get("/:organization_id", h.List)
}
