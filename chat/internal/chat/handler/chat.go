package handler

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// Chat handles chat ws functional.
//
//nolint:funlen // will be soon refactored
func (h *Handler) Chat(c *websocket.Conn) { //nolint:gocyclo // will be soon refactored
	log.Info().Str("addr", c.LocalAddr().String()).Msg("new conn")
	c.SetCloseHandler(func(code int, text string) error {
		log.Info().Int("code", code).Str("text", text).Msg("conn closed")
		return nil
	})
	defer func() {
		if err := c.Close(); err != nil {
			log.Warn().Err(err).Msg("failed to close websocket conn")
		}
	}()

	mt, msg, tokenErr := c.ReadMessage()
	if tokenErr != nil {
		log.Error().Err(tokenErr).Msg("failed to read first message")
		return
	}
	if mt != websocket.TextMessage {
		respondRaw(c, "need authorization first", errors.New("unexpected message type"))
		return
	}

	ctx, username, authErr := h.ctrl.Authorize(context.Background(), string(msg))
	if authErr != nil {
		respondRaw(c, "unauthorized", authErr)
		return
	}

	sessionID, uuidErr := uuid.Parse(c.Params("session_id"))
	if uuidErr != nil {
		respondRaw(c, "invalid session uuid", uuidErr)
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
			respondRaw(c, "failed to read query", err)
			return
		}

		if req.Command == shared.CommandCancel {
			cancel <- struct{}{}
			log.Debug().Msg("predict was canceled")
			continue
		}

		select {
		case queryID := <-hint:
			log.Debug().Msg("processing hint")
			query, err := h.ctrl.Hint(ctx, queryID, req)
			if err != nil {
				respondRaw(c, "failed to process hint", err)
				hint <- queryID
				continue
			}
			queryMsg, err := json.Marshal(query)
			if err != nil {
				respondRaw(c, err.Error(), err)
				return
			}
			respondRaw(c, string(queryMsg), nil)
			validate <- queryID
			continue
		case queryID := <-validate:
			if req.Command == shared.CommandValid {
				// for debug purposes
				log.Debug().Msg("extracted prompt is valid")
				go h.ctrl.Predict(ctx, out, cancel, queryID)
				if err := h.ctrl.UpdateStatus(ctx, queryID, shared.StatusValid); err != nil {
					respondRaw(c, "failed to update status to valid", err)
					validate <- queryID
					cancel <- struct{}{}
					continue
				}
			} else if req.Command == shared.CommandInvalid {
				log.Debug().Msg("extracted prompt is invalid")
				if err := h.ctrl.UpdateStatus(ctx, queryID, shared.StatusPending); err != nil {
					respondRaw(c, "failed to update status to pending", err)
					validate <- queryID
					continue
				}
				hint <- queryID
				continue
			}
		default:
			query, err := h.ctrl.InsertQuery(ctx, req, username, sessionID)
			if err != nil {
				respondRaw(c, err.Error(), err)
				return
			}
			queryMsg, err := json.Marshal(query)
			if err != nil {
				respondRaw(c, err.Error(), err)
				return
			}
			respondRaw(c, string(queryMsg), nil)

			validate <- query.ID
			continue
		}

		go func() {
			for {
				select {
				case chunk := <-out:
					respond(c, chunk)
					if chunk.Finish {
						log.Debug().Msg("predict finished")
						return
					}
				default:
					time.Sleep(time.Second * 1)
					log.Debug().Msg("waiting for messages...")
				}
			}
		}()
	}
}
