package organization

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	authmw "github.com/yogenyslav/ldt-2024/admin/internal/auth/middleware"
)

type organizationHandler interface {
	InsertOne(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	ImportData(c *fiber.Ctx) error
}

// SetupOrganizationRoutes устанавливает маршруты для организаций.
func SetupOrganizationRoutes(app *fiber.App, h organizationHandler, kc *gocloak.GoCloak, realm, cipher string, repo authmw.UserOrganizationRepo) {
	g := app.Group("/admin/organization")
	g.Use(authmw.JWT(kc, realm, cipher, repo))

	g.Post("/", h.InsertOne)
	g.Get("/", h.FindOne)
	g.Put("/", h.UpdateOne)
	g.Post("/import", h.ImportData)
}
