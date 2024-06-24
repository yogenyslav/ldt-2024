package favorite

import (
	"github.com/gofiber/fiber/v2"
)

type favoriteHandler interface {
	InsertOne(c *fiber.Ctx) error
	List(c *fiber.Ctx) error
	FindOne(c *fiber.Ctx) error
	UpdateOne(c *fiber.Ctx) error
	DeleteOne(c *fiber.Ctx) error
}

// SetupFavoriteRoutes настраивает маршруты для избранных предиктов.
func SetupFavoriteRoutes(app fiber.Router, h favoriteHandler) {
	g := app.Group("/favorite")

	g.Post("/", h.InsertOne)
	g.Get("/list", h.List)
	g.Get("/:id", h.FindOne)
	g.Put("/", h.UpdateOne)
	g.Delete("/:id", h.DeleteOne)
}
