package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
	chatresp "github.com/yogenyslav/ldt-2024/chat/pkg/chat_response"
)

// GetOrganization получить выбранную организацию.
func (h *Handler) GetOrganization(c *websocket.Conn) (string, error) {
	mt, msg, err := c.ReadMessage()
	if err != nil {
		log.Error().Err(err).Msg("failed to read first message")
		return "", err
	}
	if mt != websocket.TextMessage {
		e := errors.New("unexpected message type")
		chatresp.RespondError(c, "need organization choice", e)
		return "", e
	}

	return fmt.Sprintf("organization-%s", string(msg)), nil
}
