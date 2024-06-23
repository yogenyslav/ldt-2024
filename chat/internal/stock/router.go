package stock

import (
	"github.com/Nerzal/gocloak/v13"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/auth/middleware"
)

type stockHandler interface {
	UniqueCodes(c *fiber.Ctx) error
}

// SetupStockRoutes маппинг запросов к товарам.
func SetupStockRoutes(app *fiber.App, h stockHandler, kc *gocloak.GoCloak, realm, cipher string) {
	g := app.Group("/chat/stock")
	g.Use(middleware.JWT(kc, realm, cipher))

	g.Get("/unique_codes/:organization_id", h.UniqueCodes)
}
