package user

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	authmw "github.com/yogenyslav/ldt-2024/admin/internal/auth/middleware"
)

type userHandler interface {
	NewUser(c *fiber.Ctx) error
	InsertOrganization(c *fiber.Ctx) error
	DeleteOrganization(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
}

// SetupUserRoutes инициализирует роуты для работы с пользователями.
func SetupUserRoutes(app *fiber.App, h userHandler, kc *gocloak.GoCloak, realm, cipher string) {
	g := app.Group("/admin/user")
	g.Use(authmw.JWT(kc, realm, cipher))

	g.Post("/", h.NewUser)
	g.Post("/organization", h.InsertOrganization)
	g.Delete("/:username", h.DeleteOrganization)
	g.Get("/:organization", h.List)
}
