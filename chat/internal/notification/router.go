package notification

import (
	"github.com/gofiber/fiber/v2"
)

type notificationHandler interface {
	Switch(c *fiber.Ctx) error
	Check(c *fiber.Ctx) error
}

// SetupNotificationRoutes маппинг запросов к уведомлениям.
func SetupNotificationRoutes(app fiber.Router, h notificationHandler) {
	g := app.Group("/notification")

	g.Post("/switch", h.Switch)
	g.Get("/check/:organization_id", h.Check)
}
