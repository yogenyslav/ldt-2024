package stock

import (
	"github.com/gofiber/fiber/v2"
)

type stockHandler interface {
	UniqueCodes(c *fiber.Ctx) error
}

// SetupStockRoutes маппинг запросов к товарам.
func SetupStockRoutes(app fiber.Router, h stockHandler) {
	g := app.Group("/stock")

	g.Get("/unique_codes/:organization_id", h.UniqueCodes)
}
