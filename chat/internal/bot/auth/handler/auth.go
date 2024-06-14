package handler

import (
	"fmt"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	tele "gopkg.in/telebot.v3"
)

// Auth хендлер для команды /auth.
func (h *Handler) Auth(c tele.Context) error {
	_, span := h.tracer.Start(pkg.GetTraceCtx(c), "Handler.Auth.Send")
	defer span.End()

	url := fmt.Sprintf("Авторизация по http://localhost:5173?tg_id=%d", c.Sender().ID)

	return c.Send(url, &tele.SendOptions{
		DisableWebPagePreview: true,
	})
}
