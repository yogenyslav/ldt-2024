package auth

import (
	"github.com/gofiber/fiber/v2"
)

type authHandler interface {
	Login(c *fiber.Ctx) error
}

// SetupAuthRoutes устанавливает маршруты для авторизации.
func SetupAuthRoutes(app *fiber.App, h authHandler) {
	g := app.Group("/admin/auth")

	g.Post("/login", h.Login)
}
