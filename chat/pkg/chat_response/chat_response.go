package chatresp

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/rs/zerolog/log"
	"github.com/yogenyslav/ldt-2024/chat/internal/shared"
)

// Response модель ответов чата.
type Response struct {
	Err      string           `json:"err,omitempty"`
	Data     any              `json:"data,omitempty"`
	Msg      string           `json:"msg,omitempty"`
	Chunk    bool             `json:"chunk"`
	DataType shared.QueryType `json:"data_type,omitempty"`
	Finish   bool             `json:"finish"`
}

func RespondError(c *websocket.Conn, msg string, err error) {
	resp := Response{
		Msg: msg,
		Err: "ok",
	}
	if err != nil {
		resp.Finish = true
		resp.Err = err.Error()
	}
	if e := c.WriteJSON(resp); e != nil {
		log.Warn().Err(e).Msg("failed to write response")
	}
}

func RespondData(c *websocket.Conn, data any) {
	resp := Response{
		Data: data,
	}
	if e := c.WriteJSON(resp); e != nil {
		log.Warn().Err(e).Msg("failed to write response")
	}
}

func Respond(c *websocket.Conn, resp Response) {
	if e := c.WriteJSON(resp); e != nil {
		log.Warn().Err(e).Msg("failed to write response")
	}
}
