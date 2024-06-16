package handler

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
	tele "gopkg.in/telebot.v3"
)

// Chat обрабатывает сообщения от пользователя.
func (h *Handler) Chat(c tele.Context) error {
	ctx, span := h.tracer.Start(pkg.GetTraceCtx(c), "Handler.Chat")
	defer span.End()

	userID := c.Sender().ID
	username, ok := c.Get(shared.UsernameKey).(string)
	if !ok {
		return errors.New("invalid username")
	}
	token, ok := c.Get(shared.TokenKey).(string)
	if !ok {
		return errors.New("invalid token")
	}
	ctx = pkg.PushToken(ctx, token)

	state, err := h.machine.GetState(ctx, userID)
	if err != nil {
		log.Error().Err(err).Msg("failed to get state")
		return err
	}

	sessionID, err := h.machine.GetSessionID(ctx, userID)
	if errors.Is(err, redis.Nil) {
		sessionID = uuid.New()
		if err = h.sc.NewSession(ctx, sessionID, username, true, userID); err != nil {
			log.Error().Err(err).Msg("failed to create session")
			return err
		}
		if err = h.machine.SetSessionID(ctx, userID, sessionID); err != nil {
			log.Error().Err(err).Msg("failed to store sessionID")
			return err
		}
	} else if err != nil {
		log.Error().Err(err).Msg("failed to get session")
		return err
	}
	defer func() {
		if err := h.cc.SessionCleanup(ctx, sessionID); err != nil {
			log.Error().Err(err).Msg("failed to cleanup session")
		}
	}()

	var (
		out chan chatresp.Response
	)

	switch state {
	case shared.StatePending:
		req := model.QueryCreateReq{
			Prompt: c.Text(),
		}
		query, err := h.cc.InsertQuery(ctx, req, username, sessionID)
		if err != nil {
			log.Error().Err(err).Msg("failed to insert query")
			return err
		}
		if err := h.sendQuery(c, query); err != nil {
			return err
		}
		if err = h.machine.SetState(ctx, userID, shared.StateValidate); err != nil {
			log.Error().Err(err).Msg("failed to set state")
			return err
		}
		if err = h.machine.SetValidateQuery(ctx, userID, query.ID); err != nil {
			log.Error().Err(err).Msg("failed to set validate query")
			return err
		}
		return nil
	case shared.StateValidate:
		log.Debug().Str("validation", c.Text()).Msg("got validation")
		queryID, err := h.machine.GetValidateQuery(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch validate query")
			return err
		}
		if c.Text() == string(shared.CommandValid) {
			defer func() {
				if err := h.machine.SetState(ctx, userID, shared.StatePending); err != nil {
					log.Error().Err(err).Msg("failed to set state")
				}
			}()
			if err := h.cc.UpdateStatus(ctx, queryID, shared.StatusValid); err != nil {
				log.Error().Err(err).Msg("failed to update status")
				return err
			}
			out = make(chan chatresp.Response)
			cancel := make(chan struct{}, 1)
			defer close(cancel)
			defer close(out)

			go h.cc.Predict(ctx, out, cancel, queryID)
		} else {
			if err := h.cc.UpdateStatus(ctx, queryID, shared.StatusInvalid); err != nil {
				log.Error().Err(err).Msg("failed to update status")
				return err
			}
			if err := h.machine.SetHintQuery(ctx, userID, queryID); err != nil {
				log.Error().Err(err).Msg("failed to set hint query")
				return err
			}
			if err := h.machine.SetState(ctx, userID, shared.StateHint); err != nil {
				log.Error().Err(err).Msg("failed to set state")
				return err
			}
		}
	case shared.StateHint:
		log.Debug().Str("hint", c.Text()).Msg("got hint")
		queryID, err := h.machine.GetHintQuery(ctx, userID)
		if err != nil {
			log.Error().Err(err).Msg("failed to fetch hint query")
			return err
		}
		query, err := h.cc.Hint(ctx, queryID, model.QueryCreateReq{Prompt: c.Text()})
		if err != nil {
			log.Error().Err(err).Msg("failed to process hint")
			return err
		}
		if err := h.sendQuery(c, query); err != nil {
			return err
		}
		if err := h.machine.SetState(ctx, userID, shared.StateValidate); err != nil {
			log.Error().Err(err).Msg("failed to set state")
			return err
		}
	default:
		return c.Send("Не могу выполнить эту команду сейчас")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go h.processChunk(&wg, c, out)
	wg.Wait()

	return nil
}
