package auth

import (
	tele "gopkg.in/telebot.v3"
)

type authHandler interface {
	Auth(c tele.Context) error
}

// SetupAuthRoutes sets up the routes for the auth handler.
func SetupAuthRoutes(b *tele.Bot, h authHandler) {
	b.Handle("/bot/auth", h.Auth)
}
