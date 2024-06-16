package handler

import (
	"context"
	"errors"

	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
)

// Authorize проверяет авторизацию пользователя.
func (h *Handler) Authorize(c *websocket.Conn) (context.Context, string, error) {
	mt, msg, tokenErr := c.ReadMessage()
	if tokenErr != nil {
		log.Error().Err(tokenErr).Msg("failed to read first message")
		return nil, "", tokenErr
	}
	if mt != websocket.TextMessage {
		chatresp.RespondError(c, "need authorization first", errors.New("unexpected message type"))
		return nil, "", errors.New("unexpected message type")
	}

	ctx, username, authErr := h.ctrl.Authorize(context.Background(), string(msg))
	if authErr != nil {
		chatresp.RespondError(c, "unauthorized", authErr)
		return nil, "", authErr
	}
	return ctx, username, nil
}
