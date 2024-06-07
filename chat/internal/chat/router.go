package chat

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/middleware"
)

type chatHandler interface {
	Chat(c *websocket.Conn)
}

// SetupChatRoutes maps the chat routes to the chat handler.
func SetupChatRoutes(app *fiber.App, h chatHandler, cfg websocket.Config) {
	g := app.Group("/chat")
	g.Use(middleware.WsProtocolUpgrade())

	g.Get("/:session_id", websocket.New(h.Chat, cfg))
}
