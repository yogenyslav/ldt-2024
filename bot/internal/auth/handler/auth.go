package handler

import (
	"github.com/yogenyslav/ldt-2024/bot/pkg"
	tele "gopkg.in/telebot.v3"
)

// Auth is the handler for the auth command.
func (h *Handler) Auth(c tele.Context) error {
	ctx, span := h.tracer.Start(pkg.GetTraceCtx(c), "Handler.Auth.Send")
	defer span.End()

	_ = ctx

	return c.Send("Авторизация по <a href=\"https://ya.ru/\">ссылке</a>", &tele.SendOptions{
		DisableWebPagePreview: true,
	})
}
