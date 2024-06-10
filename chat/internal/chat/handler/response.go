package handler

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
)

// Response is a struct for websocket error responses.
type Response struct {
	Err error  `json:"err,omitempty"`
	Msg string `json:"msg"`
}

func respondRaw(c *websocket.Conn, msg string, err error) {
	if e := c.WriteJSON(Response{
		Msg: msg,
		Err: err,
	}); e != nil {
		log.Warn().Err(e).Msg("failed to write response")
	}
}

func respond(c *websocket.Conn, resp Response) {
	if e := c.WriteJSON(resp); e != nil {
		log.Warn().Err(e).Msg("failed to write response")
	}
}
