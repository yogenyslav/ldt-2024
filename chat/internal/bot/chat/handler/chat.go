package handler

import (
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	tele "gopkg.in/telebot.v3"
)

// Chat обрабатывает сообщения от пользователя.
func (h *Handler) Chat(c tele.Context) error {
	ctx, span := h.tracer.Start(pkg.GetTraceCtx(c), "Handler.Chat")
	defer span.End()

	username := c.Get(shared.UsernameKey).(string)

	state, err := h.machine.GetState(ctx, c.Sender().ID)
	if err != nil {
		return err
	}

	sessionID, err := h.machine.GetSessionID(ctx, c.Sender().ID)
	if err != nil {
		return err
	}

	if state == shared.StatePending {
		req := model.QueryCreateReq{
			Prompt: c.Text(),
		}
		query, err := h.cc.InsertQuery(ctx, req, username, sessionID)
		if err != nil {
			return err
		}
		if err := c.Send(query); err != nil {
			return err
		}
	}
	return nil
}
