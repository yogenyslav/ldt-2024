package handler

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

func closeHandler(code int, text string) error {
	log.Info().Int("code", code).Str("text", text).Msg("close handler")
	return nil
}

// Chat обрабатывает вебсокет соединение.
func (h *Handler) Chat(c *websocket.Conn) {
	log.Info().Str("addr", c.LocalAddr().String()).Msg("new conn")
	c.SetCloseHandler(closeHandler)
	defer func() {
		if err := c.Close(); err != nil {
			log.Warn().Err(err).Msg("failed to close websocket conn")
		}
		log.Info().Msg("conn closed")
	}()

	ctx, username, authErr := h.Authorize(c)
	if authErr != nil {
		return
	}
	sessionID, uuidErr := uuid.Parse(c.Params("session_id"))
	if uuidErr != nil {
		respondError(c, "invalid session uuid", uuidErr)
		return
	}
	defer func() {
		if err := h.ctrl.SessionCleanup(ctx, sessionID); err != nil {
			log.Error().Err(err).Msg("failed to cleanup session")
		}
	}()

	var (
		validate = make(chan int64, 1)
		out      = make(chan Response)
		hint     = make(chan int64, 1)
		cancel   = make(chan struct{}, 1)
	)

	for {
		var req model.QueryCreateReq
		if err := c.ReadJSON(&req); err != nil {
			respondError(c, "failed to read query", err)
			return
		}
		if req.Command == shared.CommandCancel {
			cancel <- struct{}{}
			continue
		}

		select {
		case queryID := <-hint:
			query, err := h.ctrl.Hint(ctx, queryID, req)
			if err != nil {
				respondError(c, "failed to process hint", err)
				hint <- queryID
				continue
			}
			respondData(c, query)
			validate <- queryID
			continue
		case queryID := <-validate:
			if req.Command == shared.CommandValid {
				go h.ctrl.Predict(ctx, out, cancel, queryID)
				if err := h.ctrl.UpdateStatus(ctx, queryID, shared.StatusValid); err != nil {
					respondError(c, "failed to update status to valid", err)
					validate <- queryID
					cancel <- struct{}{}
					continue
				}
			} else if req.Command == shared.CommandInvalid {
				if err := h.ctrl.UpdateStatus(ctx, queryID, shared.StatusInvalid); err != nil {
					respondError(c, "failed to update status to invalid", err)
					validate <- queryID
					continue
				}
				hint <- queryID
				continue
			}
		default:
			query, err := h.ctrl.InsertQuery(ctx, req, username, sessionID)
			if err != nil {
				respondError(c, err.Error(), err)
				return
			}
			respondData(c, query)
			validate <- query.ID
			continue
		}
		go h.ProcessChunks(out, c)
	}
}
