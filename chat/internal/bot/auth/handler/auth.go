package handler

import (
	"github.com/yogenyslav/ldt-2024/bot/pkg"
	tele "gopkg.in/telebot.v3"
)

// Auth хендлер для команды /auth.
func (h *Handler) Auth(c tele.Context) error {
	_, span := h.tracer.Start(pkg.GetTraceCtx(c), "Handler.Auth.Send")
	defer span.End()

	return c.Send("Авторизация по <a href=\"https://ya.ru/\">ссылке</a>", &tele.SendOptions{
		DisableWebPagePreview: true,
	})
}
