package handler

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/chat/model"
	tele "gopkg.in/telebot.v3"
)

func (h *Handler) sendQuery(c tele.Context, query model.QueryDto) error {
	if err := c.Send(fmt.Sprintf("%v", query)); err != nil {
		log.Error().Err(err).Msg("failed to send query meta")
		return err
	}
	return nil
}
