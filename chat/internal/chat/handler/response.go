package handler

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
)

// ErrorResponse is a struct for websocket error responses.
type ErrorResponse struct {
	Err error  `json:"err"`
	Msg string `json:"msg"`
}

func writeError(c *websocket.Conn, msg string, err error) {
	if e := c.WriteJSON(ErrorResponse{
		Msg: msg,
		Err: err,
	}); e != nil {
		log.Warn().Err(e).Msg("failed to write error response")
	}
	return
}
