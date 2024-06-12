package handler

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
)

// Response модель ответов чата.
type Response struct {
	Err    string `json:"err,omitempty"`
	Data   any    `json:"data,omitempty"`
	Msg    string `json:"msg,omitempty"`
	Chunk  bool   `json:"chunk"`
	Finish bool   `json:"finish"`
}

func respondError(c *websocket.Conn, msg string, err error) {
	resp := Response{
		Msg: msg,
		Err: err.Error(),
	}
	if err != nil {
		resp.Finish = true
	}
	if e := c.WriteJSON(resp); e != nil {
		log.Warn().Err(e).Msg("failed to write response")
	}
}

func respondData(c *websocket.Conn, data any) {
	resp := Response{
		Data: data,
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
