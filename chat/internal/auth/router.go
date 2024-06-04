package auth

import (
	"github.com/gofiber/fiber/v2"
)

type authHandler interface {
	Login(c *fiber.Ctx) error
}

func SetupAuthRoutes(app *fiber.App, h authHandler) {
	g := app.Group("/auth")

	g.Post("/login", h.Login)
}
