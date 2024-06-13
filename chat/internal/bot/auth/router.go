package auth

import (
	tele "gopkg.in/telebot.v3"
)

type authHandler interface {
	Auth(c tele.Context) error
}

// SetupAuthRoutes маппинг путей авторизации.
func SetupAuthRoutes(b *tele.Bot, h authHandler) {
	b.Handle("/auth", h.Auth)
}
