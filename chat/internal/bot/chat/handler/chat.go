package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
	"github.com/yogenyslav/ldt-2024/chat/pkg"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
	tele "gopkg.in/telebot.v3"
	"sync"
	"time"
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
		if err = c.Send(fmt.Sprintf("%v", query)); err != nil {
			log.Error().Err(err).Msg("failed to send query meta")
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
		if c.Text() == string(shared.CommandValid) {
			defer func() {
				if err := h.machine.SetState(ctx, userID, shared.StatePending); err != nil {
					log.Error().Err(err).Msg("failed to set state")
				}
			}()
			queryID, err := h.machine.GetValidateQuery(ctx, userID)
			if err != nil {
				log.Error().Err(err).Msg("failed to fetch validate query")
				return err
			}
			out = make(chan chatresp.Response)
			cancel := make(chan struct{}, 1)
			defer close(cancel)
			defer close(out)

			go h.cc.Predict(ctx, out, cancel, queryID)
		}
	default:
		return c.Send("Не могу выполнить эту команду сейчас")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-time.After(time.Minute):
				if err := c.Send("response timeout exceeded"); err != nil {
					log.Error().Err(err).Msg("failed to send response")
				}
				return
			case chunk := <-out:
				if chunk.Err != "" {
					log.Error().Err(errors.New(chunk.Err)).Msg("failed to predict")
					if err := c.Send("Не удалось получить данные"); err != nil {
						log.Error().Err(err).Msg("failed to send response")
						continue
					}
					return
				}
				msg, err := json.Marshal(chunk)
				if err != nil {
					log.Error().Err(err).Msg("failed to marshal response")
					continue
				}
				if err := c.Send(string(msg)); err != nil {
					if errors.Is(err, tele.ErrTooLongMessage) {
						_ = c.Send(string(msg)[:300])
					}
					log.Error().Err(err).Msg("failed to send chunk")
					continue
				}
				if chunk.Finish {
					return
				}
			}
		}
	}()
	wg.Wait()

	return nil
}
