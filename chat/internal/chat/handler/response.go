package handler

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
)

// Response is a struct for websocket error responses.
type Response struct {
	Err    error  `json:"err,omitempty"`
	Msg    string `json:"msg"`
	Finish bool   `json:"finish"`
}

func respondRaw(c *websocket.Conn, msg string, err error) {
	resp := Response{
		Msg: msg,
		Err: err,
	}
	if err != nil {
		resp.Finish = true
	}
	if e := c.WriteJSON(resp); e != nil {
		log.Warn().Err(e).Msg("failed to write response")
	}
}

func respond(c *websocket.Conn, resp Response) {
	if e := c.WriteJSON(resp); e != nil {
		log.Warn().Err(e).Msg("failed to write response")
	}
}
