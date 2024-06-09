package handler

import (
	"context"
	"errors"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
)

// Chat handles chat ws functional.
func (h *Handler) Chat(c *websocket.Conn) {
	ctx := context.Background()

	log.Info().Str("addr", c.LocalAddr().String()).Msg("new conn")
	c.SetCloseHandler(func(code int, text string) error {
		log.Info().Int("code", code).Str("text", text).Msg("conn closed")
		return nil
	})
	defer func() {
		if err := c.Close(); err != nil {
			log.Warn().Err(err).Msg("failed to close websocket conn")
			return
		}
	}()

	mt, msg, tokenErr := c.ReadMessage()
	if tokenErr != nil {
		log.Error().Err(tokenErr).Msg("failed to read first message")
		return
	}
	if mt != websocket.TextMessage {
		writeError(c, "need authorization first", errors.New("unexpected message type"))
		return
	}

	username, authErr := h.ctrl.Authorize(ctx, string(msg))
	if authErr != nil {
		writeError(c, "unauthorized", authErr)
		return
	}

	sessionID, uuidErr := uuid.Parse(c.Params("session_id"))
	if uuidErr != nil {
		writeError(c, "invalid session uuid", uuidErr)
		return
	}

	var (
		validate = make(chan struct{}, 1)
	)

	for {
		var req model.QueryCreateReq
		if err := c.ReadJSON(&req); err != nil {
			writeError(c, "failed to read query", err)
			return
		}

		select {
		case <-validate:
			if req.Command == "valid" {
				log.Debug().Msg("extracted prompt is valid")
			}
		default:
			validate <- struct{}{}
			if err := h.ctrl.InsertQuery(ctx, req, username, sessionID); err != nil {
				writeError(c, err.Error(), err)
				return
			}
		}
	}
}
