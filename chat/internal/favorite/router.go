package favorite

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/middleware"
)

type favoriteHandler interface {
	InsertOne(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
}

// SetupFavoriteRoutes настраивает маршруты для избранных предиктов.
func SetupFavoriteRoutes(app *fiber.App, h favoriteHandler, kc *gocloak.GoCloak, realm, cipher string) {
	g := app.Group("/chat/favorite")
	g.Use(middleware.JWT(kc, realm, cipher))

	g.Post("/", h.InsertOne)
	g.Get("/list", h.List)
	g.Get("/:id", h.FindOne)
	g.Put("", h.UpdateOne)
	g.Delete("/:id", h.DeleteOne)
}
