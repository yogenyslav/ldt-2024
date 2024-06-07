package auth

import (
	tele "gopkg.in/telebot.v3"
)

type authHandler interface {
	Auth(c tele.Context) error
}

func SetupUserRoutes(b *tele.Bot, h authHandler) {
	b.Handle("/auth", h.Auth)
}
