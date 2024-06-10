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
//
//nolint:funlen // will be soon refactored
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
		respond(c, "need authorization first", errors.New("unexpected message type"))
		return
	}

	username, authErr := h.ctrl.Authorize(ctx, string(msg))
	if authErr != nil {
		respond(c, "unauthorized", authErr)
		return
	}

	sessionID, uuidErr := uuid.Parse(c.Params("session_id"))
	if uuidErr != nil {
		respond(c, "invalid session uuid", uuidErr)
		return
	}

	var (
		validate = make(chan int64, 1)
		out      = make(chan Response)
		cancel   = make(chan struct{}, 1)
	)
	defer close(out)

	for {
		var req model.QueryCreateReq
		if err := c.ReadJSON(&req); err != nil {
			respond(c, "failed to read query", err)
			return
		}

		if req.Command == "cancel" {
			cancel <- struct{}{}
			log.Debug().Msg("predict was canceled")
			continue
		}

		select {
		case queryID := <-validate:
			if req.Command == "valid" {
				// for debug purposes
				log.Debug().Msg("extracted prompt is valid")
				go h.ctrl.Predict(ctx, out, cancel, queryID)
			} else {
				log.Debug().Msg("extracted prompt is invalid")
			}
		default:
			queryID, err := h.ctrl.InsertQuery(ctx, req, username, sessionID)
			if err != nil {
				respond(c, err.Error(), err)
				return
			}
			validate <- queryID
		}
	}
}
