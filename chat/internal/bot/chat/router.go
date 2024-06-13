package chat

import (
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/chat/middleware"
	"github.com/yogenyslav/ldt-2024/chat/internal/bot/state"
	tele "gopkg.in/telebot.v3"
)

type chatHandler interface {
	Chat(c tele.Context) error
}

// SetupChatRoutes маппинг путей чата.
func SetupChatRoutes(g *tele.Group, h chatHandler, machine *state.Machine) {
	g.Use(middleware.State(machine))
	g.Handle(tele.OnText, h.Chat)
}
