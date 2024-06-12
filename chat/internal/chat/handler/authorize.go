package handler

import (
	"context"
	"errors"

	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
)

// Authorize проверяет авторизацию пользователя.
func (h *Handler) Authorize(c *websocket.Conn) (context.Context, string, error) {
	mt, msg, tokenErr := c.ReadMessage()
	if tokenErr != nil {
		log.Error().Err(tokenErr).Msg("failed to read first message")
		return nil, "", tokenErr
	}
	if mt != websocket.TextMessage {
		respondError(c, "need authorization first", errors.New("unexpected message type"))
		return nil, "", errors.New("unexpected message type")
	}

	ctx, username, authErr := h.ctrl.Authorize(context.Background(), string(msg))
	if authErr != nil {
		respondError(c, "unauthorized", authErr)
		return nil, "", authErr
	}
	return ctx, username, nil
}
