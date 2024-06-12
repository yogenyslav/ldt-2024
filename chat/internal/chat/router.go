package chat

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/middleware"
)

type chatHandler interface {
	Chat(c *websocket.Conn)
}

// SetupChatRoutes маппит пути для чата.
func SetupChatRoutes(app *fiber.App, h chatHandler, cfg websocket.Config) {
	g := app.Group("/chat/ws")
	g.Use(middleware.WsProtocolUpgrade())

	g.Get("/:session_id", websocket.New(h.Chat, cfg))
}
